package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	var depts []string
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			must(err)
		}
		depts = append(depts, strings.TrimSpace(s))
	}

	var courses []Course
	for _, dept := range depts {
		courses = append(courses, getCourses(dept)...)
	}
	dump(courses)
}
