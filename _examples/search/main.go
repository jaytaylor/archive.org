package main

import (
	"fmt"

	"jaytaylor.com/archive.org"
)

func main() {
	u := "http://blog.sendhub.com/post/16800984141/switching-to-heroku-a-django-app-story"

	hits, err := archiveorg.Search(u, archiveorg.DefaultRequestTimeout)
	if err != nil {
		panic(fmt.Errorf("Search error: %s", err))
	}
	fmt.Printf("num: %v\n", len(hits))
	for _, hit := range hits {
		fmt.Printf("hit: %+v\n", hit)
	}
}
