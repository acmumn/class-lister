package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var courseNumRegex = regexp.MustCompile("[0-9]{3,4}")
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

func getCourses(dept string) []Course {
	courseResp, err := http.Get("https://www.aem.umn.edu/cgi-bin/courses/noauth/class-schedule?current=" + dept)
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
			courseNumMatches := courseNumRegex.FindStringSubmatch(children.Find("b").Text())
			if len(courseNumMatches) == 0 {
				log.Println("Warning: Couldn't find course number")
				html, _ := e.Html()
				log.Println("    ", dept)
				log.Println("    ", html)
			}
			courseNum, err := strconv.Atoi(courseNumMatches[0])
			must(err)
			courses = append(courses, Course{
				Number:   courseNum,
				Lectures: nil,
			})
		} else if childCount == 5 {
			classType := e.Find(":nth-child(2)").Text()
			if classType != "Lecture" {
				return
			}
			timeLocStr := strings.TrimSpace(unfuck(e.Find(":nth-child(3)")))
			if timeLocStr == "" || strings.Contains(timeLocStr, "(online)") {
				return
			}
			timeLocData := timeLocRegex.FindStringSubmatch(timeLocStr)
			if timeLocData == nil {
				log.Printf("Warning: Couldn't parse %#v", timeLocStr)
				html, _ := e.Html()
				log.Println("    ", dept)
				log.Println("    ", html)
				return
			}
			startTime := timeLocData[1]
			days := timeLocData[2]
			room := timeLocData[3]
			instructor := unfuck(e.Find(":nth-child(4)"))

			if len(courses) == 0 {
				log.Println("Something screwy is going on with the page...")
				log.Printf("Check on the dept %#v", dept)
				return
			}

			c := &courses[len(courses)-1]
			c.Lectures = append(c.Lectures, Lecture{
				Time:      startTime,
				Days:      days,
				Room:      room,
				Professor: instructor,
			})
		} else {
			log.Println("Warning: Unrecognized row")
			html, _ := e.Html()
			log.Println("    ", dept)
			log.Println("    ", html)
		}
	})
	return courses
}
