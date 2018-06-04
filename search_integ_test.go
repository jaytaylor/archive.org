// +build integration

package archiveorg

import (
	"testing"
)

func TestSearch(t *testing.T) {
	u := "http://blog.sendhub.com/post/16800984141/switching-to-heroku-a-django-app-story"

	hits, err := Search(u)
	if err != nil {
		t.Fatalf("Search error: %s", err)
	}
	t.Logf("hits: %+v", hits)
}
