package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unfuck(s *goquery.Selection) string {
	t := s.Text()
	t = strings.Replace(t, "\u00a0", " ", -1)
	t = strings.Replace(t, "\u2011", "-", -1)
	return strings.TrimSpace(t)
}
