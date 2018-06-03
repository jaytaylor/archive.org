package main

import (
	"fmt"

	"github.com/jaytaylor/archive.org"
)

var captureURL = "https://jaytaylor.com/"

func main() {
	archiveURL, err := archiveorg.Capture(captureURL, archiveorg.DefaultRequestTimeout)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully archived %v via archive.org: %v\n", captureURL, archiveURL)
}
