package archiveorg

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var NoContentLocationErr = errors.New("missing 'content-lcation' header") // Returned when a malformed response is returned by archive.org.

// Capture requests a
func Capture(url string, timeout time.Duration) (string, error) {
	pleaseCrawl := fmt.Sprintf("%v/save/%v", BaseURL, url)

	log.WithField("crawl-request", pleaseCrawl).Debugf("Requesting archive.org crawl")

	resp, _, err := doRequest("", pleaseCrawl, nil, timeout)
	if err != nil {
		return "", err
	}

	loc := resp.Header.Get("Content-Location")

	if loc == "" {
		return "", NoContentLocationErr
	}

	location := fmt.Sprintf("%v/%v", BaseURL, loc)

	return location, nil
}
