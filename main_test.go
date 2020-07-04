package main

import (
	"strings"
	"testing"

	"github.com/wabarc/archive.org/pkg"
)

func TestWayback(t *testing.T) {
	url := "https://www.google.com"
	links := []string{url}
	wbrc := &ia.Archiver{}
	got, _ := wbrc.Wayback(links)
	for _, dest := range got {
		if strings.Contains(dest, url) == false || strings.Contains(dest, "archive.org") == false {
			t.Error(got)
			t.Fail()
		}
	}
}
