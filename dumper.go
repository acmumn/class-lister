package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func dump(courses []Course) {
	w := csv.NewWriter(os.Stdout)
	for _, course := range courses {
		for _, lecture := range course.Lectures {
			w.Write(toRow(&course, &lecture))
		}
	}
	w.Flush()
}

func toRow(course *Course, lecture *Lecture) []string {
	return []string{
		strconv.Itoa(course.Number),
		lecture.Time,
		lecture.Days,
		lecture.Room,
		lecture.Professor,
	}
}
