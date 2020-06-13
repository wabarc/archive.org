package main

import (
	"github.com/wabarc/archive.org/pkg"
	"strings"
	"testing"
)

func TestWayback(t *testing.T) {
	url := "https://www.google.com"
	links := []string{url}
	got := ia.Wayback(links)
	s := strings.Join(got, " ")
	if strings.Contains(s, url) == false || strings.Contains(s, "archive.org") == false {
		t.Error(got)
		t.Fail()
	}
}
