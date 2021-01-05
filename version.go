package main

import "fmt"

var (
	version = "1.0.0"
	date    = "unknown"
)

func init() {
	fmt.Printf("version: %s\ndate: %s\n\n", version, date)
}
