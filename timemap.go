package archiveorg

// Time map parser as described and linked at
// http://ws-dl.blogspot.com/2013/07/2013-07-15-wayback-machine-upgrades.html.

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	MementoParseErr = errors.New("malformed input: memento parse failed")

	mementoLayout         = "Mon, 02 Jan 2006 15:04:05 MST"
	mementoOuterSplitExpr = regexp.MustCompile(`^<(.*)>; (.*),$`) //(?:([^=]+)="([^"]*)"[;,])*$`) // (?: type="([^"]*)"[;,])?(?: from="([^"]*)"[;,])?(?: datetime="([^"]*)"[;,])?$`)
	mementoTailSplitExpr  = regexp.MustCompile(` *; *`)           // ([^=]+)="([^"]*)"[;,])`)
	mementoTailParseExpr  = regexp.MustCompile(`([^=]+)="([^"]*)"`)
)

/*
<http://www.jaytaylor.com:80/>; rel="original",
<http://web.archive.org/web/timemap/link/https://jaytaylor.com>; rel="self"; type="application/link-format"; from="Sat, 31 Mar 2001 11:48:39 GMT",
<http://web.archive.org>; rel="timegate",
<http://web.archive.org/web/20010331114839/http://www.jaytaylor.com:80/>; rel="first memento"; datetime="Sat, 31 Mar 2001 11:48:39 GMT",
*/

type TimeMap struct {
	Original *Memento
	Self     *Memento
	TimeGate *Memento
	Mementos []Memento
}

type Memento struct {
	URL  string
	Rel  string
	Type *string    `json:,omitempty`
	From *time.Time `json:,omitempty`
	Time *time.Time `json:,omitempty`
}

func NewTimeMap() *TimeMap {
	timemap := &TimeMap{
		Mementos: []Memento{},
	}
	return timemap
}

func TimeMapFor(url string, timeout ...time.Duration) (*TimeMap, error) {
	if len(timeout) == 0 {
		timeout = []time.Duration{DefaultRequestTimeout}
	}

	resp, err := downloadTimeMap(url, timeout[0])
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	timemap, err := ParseTimeMap(resp.Body)
	if err != nil {
		return nil, err
	}

	return timemap, nil
}

// ParseTimeMap takes a reader and parses it as a complete TimeMap.
func ParseTimeMap(r io.Reader) (*TimeMap, error) {
	timemap := NewTimeMap()

	scanner := bufio.NewScanner(r)

	for i := 1; scanner.Scan(); i++ {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			m, err := ParseMemento(line)
			if err != nil {
				return nil, fmt.Errorf("%s: on line %v: %v", err, i, line)
			}

			switch m.Rel {
			case "memento", "first memento":
				timemap.Mementos = append(timemap.Mementos, *m)

			case "timegate":
				timemap.TimeGate = m

			case "self":
				timemap.Self = m

			case "original":
				timemap.Original = m

			case "default":
				return nil, fmt.Errorf("no handler for memento rel value %q", m.Rel)
			}
		}
	}

	return timemap, nil
}

// ParseMemento parses a line containing a Memento entry.
func ParseMemento(line string) (*Memento, error) {
	outerPieces := mementoOuterSplitExpr.FindAllStringSubmatch(line, -1)

	if len(outerPieces) == 0 {
		return nil, MementoParseErr
	}

	m := &Memento{
		URL: outerPieces[0][1],
	}

	tailPieces := mementoTailSplitExpr.Split(outerPieces[0][2], -1)

	if len(tailPieces) == 0 {
		return m, nil
	}

	for _, tailPiece := range tailPieces {
		matches := mementoTailParseExpr.FindAllStringSubmatch(tailPiece, -1)
		if len(matches) == 0 {
			log.WithField("line", line).Warnf("Unexpected input, unrecognized tail piece: %v", tailPiece)
			continue
		}

		switch matches[0][1] {
		case "rel":
			m.Rel = matches[0][2]

		case "type":
			typ := matches[0][2]
			m.Type = &typ

		case "from":
			t, err := time.Parse(mementoLayout, matches[0][2])
			if err != nil {
				log.Errorf("%s", err)
				return nil, MementoParseErr
			}
			m.From = &t

		case "datetime":
			t, err := time.Parse(mementoLayout, matches[0][2])
			if err != nil {
				log.Errorf("%s", err)
				return nil, MementoParseErr
			}
			m.Time = &t

		default:
			log.WithField("line", line).Warnf("Unexpected input, unrecognized memento field: %v", matches[0][1])
		}
	}

	return m, nil
}

func downloadTimeMap(url string, timeout time.Duration) (*http.Response, error) {
	timeMapURL := fmt.Sprintf("%v/web/timemap/link/%v", BaseURL, url)

	req, err := newRequest("", timeMapURL, nil)
	if err != nil {
		return nil, err
	}

	client := newClient(timeout)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request to %v: %s", timeMapURL, err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request to %v received non-2xx response status-code=%v", timeMapURL, resp.StatusCode)
	}

	return resp, nil
}
