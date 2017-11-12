package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var courseNumRegex = regexp.MustCompile("[0-9]{4}")
var timeLocRegex = regexp.MustCompile("([0-9:]{5} [AP]M)-[0-9:]{5} [AP]M ([^ ]+) (.+)")

type Course struct {
	Number   int
	Lectures []Lecture
}

type Lecture struct {
	Time      string
	Days      string
	Room      string
	Professor string
}

func getCourses() []Course {
	courseResp, err := http.Get("https://www.aem.umn.edu/cgi-bin/courses/noauth/class-schedule?current=CSCI")
	must(err)
	doc, err := goquery.NewDocumentFromResponse(courseResp)
	must(err)

	courses := make([]Course, 0)
	doc.Find("#content > table:nth-child(7) > tbody > tr").Each(func(n int, e *goquery.Selection) {
		if n == 0 {
			return
		}

		children := e.Children()
		childCount := children.Length()
		if childCount == 1 {
			courseNum, err := strconv.Atoi(courseNumRegex.FindStringSubmatch(children.Find("b").Text())[0])
			must(err)
			if courseNum >= 5000 {
				return
			}
			courses = append(courses, Course{
				Number:   courseNum,
				Lectures: nil,
			})
		} else if childCount == 5 {
			classType := e.Find(":nth-child(2)").Text()
			if classType != "Lecture" {
				return
			}
			timeLocStr := unfuck(e.Find(":nth-child(3)"))
			timeLocData := timeLocRegex.FindStringSubmatch(timeLocStr)
			if timeLocData == nil {
				log.Println("Warning: Couldn't parse ", timeLocStr)
				return
			}
			startTime := timeLocData[1]
			days := timeLocData[2]
			room := timeLocData[3]
			if room == "(online)" {
				return
			}
			instructor := unfuck(e.Find(":nth-child(4)"))

			c := &courses[len(courses)-1]
			c.Lectures = append(c.Lectures, Lecture{
				Time:      startTime,
				Days:      days,
				Room:      room,
				Professor: instructor,
			})
		} else {
			panic("TODO")
		}
	})
	return courses
}
