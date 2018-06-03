# archiveorg

[![Documentation](https://godoc.org/github.com/jaytaylor/archive.org?status.svg)](https://godoc.org/github.com/jaytaylor/archive.org)
[![Build Status](https://travis-ci.org/jaytaylor/archive.org.svg?branch=master)](https://travis-ci.org/jaytaylor/archiveorg)
[![Report Card](https://goreportcard.com/badge/github.com/jaytaylor/archive.org)](https://goreportcard.com/report/github.com/jaytaylor/archive.org)

### About

archive.org is a golang package for archiving web pages via [archive.org](https://web.archive.org).

Please be mindful and responsible and go easy on them, we want archive.org to last forever!

Created by [Jay Taylor](https://jaytaylor.com/).

Also see: [archive.is golang package](https://jaytaylor.com/archive.is)

### Requirements

* Go version 1.9 or newer

### Installation

```bash
go get jaytaylor.com/archive.org/...
```

### Usage

#### Command-line programs

##### `archive.org <url>`

Archive a fresh new copy of an HTML page

##### `archive.org-snapshots <url>`

Search for existing page snapshots

#### Go package interfaces

##### Search for Existing Snapshots

[capture.go](_examples/capture/capture.go):

```go
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

// Output:
//
// Successfully archived https://jaytaylor.com/ via archive.org: https://archive.is/i2PiW
```


[search.go](_examples/search/search.go):

```go
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

// Output:
//
// num: 3
// hit: {URL:https://web.archive.org/web/20160304012638/http://blog.sendhub.com/post/16800984141/switching-to-heroku-a-django-app-story Reason:webwidecrawlhackernews00000hackernews StatusCode:301 Timestamp:2016-03-04 01:26:38 +0000 UTC}
// hit: {URL:https://web.archive.org/web/20120202233158/http://blog.sendhub.com/post/16800984141/switching-to-heroku-a-django-app-story Reason:alexacrawls StatusCode:200 Timestamp:2012-02-02 23:31:58 +0000 UTC}
// hit: {URL:https://web.archive.org/web/20120202201233/http://blog.sendhub.com/post/16800984141/switching-to-heroku-a-django-app-story Reason:alexacrawls StatusCode:200 Timestamp:2012-02-02 20:12:33 +0000 UTC}
```

### Running the test suite

    go test ./...

#### License

Permissive MIT license, see the [LICENSE](LICENSE) file for more information.
