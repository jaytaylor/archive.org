package archiveorg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const timestampLayout = "20060102150405"

var (
	BaseURL               = "https://web.archive.org"                                                                                                  // Overrideable default package value.
	HTTPHost              = "archive.org"                                                                                                              // Overrideable default package value.
	UserAgent             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36" // Overrideable default package value.
	DefaultRequestTimeout = 10 * time.Second                                                                                                           // Overrideable default package value.
)

// Snapshot represents an instance of a URL page snapshot on archive.is.
type Snapshot struct {
	URL        string
	Reason     string
	StatusCode int
	Timestamp  time.Time
}

// Search for URL snapshots.
func Search(u string, timeout ...time.Duration) ([]Snapshot, error) {
	if len(timeout) == 0 {
		timeout = []time.Duration{DefaultRequestTimeout}
	}

	sl, err := sparklineFor(u, timeout[0])
	if err != nil {
		return nil, err
	}

	points, err := sl.captures()
	if err != nil {
		return nil, err
	}

	snaps := []Snapshot{}

	for _, point := range points {
		for i := 0; i < point.Count; i++ {
			if i >= len(point.Timestamps) {
				// Invalid offset, skip entry.
				continue
				log.WithField("url", u).Warn("Skipping point with missing timestamp")
			}

			snap := Snapshot{
				URL: fmt.Sprintf("%v/web/%v/%v", BaseURL, point.Timestamps[i], u),
			}

			if ts, err := time.Parse(timestampLayout, fmt.Sprint(point.Timestamps[i])); err != nil {
				log.WithField("url", u).WithField("timestamp", point.Timestamps[i]).Errorf("Skipping point after failed timestamp parse")
				continue
			} else {
				snap.Timestamp = ts
			}

			if i < len(point.Whys) {
				snap.Reason = strings.Join(point.Whys[i], "")
			}
			if i < len(point.StatusCodes) {
				if sc, err := strconv.Atoi(fmt.Sprint(point.StatusCodes[i])); err == nil {
					snap.StatusCode = sc
				}
			}
			snaps = append(snaps, snap)
		}
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].Timestamp.After(snaps[j].Timestamp)
	})

	return snaps, nil
}

type calendarPoint struct {
	Count       int           `json:"cnt"`
	Whys        [][]string    `json:"why"`
	StatusCodes []interface{} `json:"st"`
	Timestamps  []int64       `json:"ts"`
}

func (point calendarPoint) isEmpty() bool {
	empty := point.Count == 0 && len(point.Whys) == 0 && len(point.StatusCodes) == 0 && len(point.Timestamps) == 0
	return empty
}

func calendarFor(u string, timeout time.Duration) ([]calendarPoint, error) {
	sl, err := sparklineFor(u, timeout)
	if err != nil {
		return nil, err
	}

	captures, err := sl.captures()

	return captures, nil
}

type sparkline struct {
	FirstTs string        `json:"first_ts"`
	LastTs  string        `json:"last_ts"`
	Years   map[int][]int `json:"years"`
	safeURL string        `json:"-"`
	timeout time.Duration `json:"-"`
}

func (sl *sparkline) captures() ([]calendarPoint, error) {
	var (
		points = []calendarPoint{}
	)

	for year, monthCounts := range sl.Years {
		for _, count := range monthCounts {
			// Capture each year with a non-empty crawl count month.
			if count > 0 {
				var (
					queryURL = fmt.Sprintf("%v/__wb/calendarcaptures?url=%v&selected_year=%v", BaseURL, sl.safeURL, year)
					captures = [][][]*calendarPoint{}
				)

				if _, err := simpleHTTPJSON(queryURL, &captures, sl.timeout); err != nil {
					return nil, err
				}

				for _, pointlessArray := range captures {
					// NB: Not clear why the API always encloses contents in an
					// array of size 1.
					if len(pointlessArray) > 0 {
						for _, point := range pointlessArray[0] {
							if point != nil && !point.isEmpty() {
								points = append(points, *point)
							}
						}
					}
				}

				// Skip the rest of the year since we already have it now.
				break
			}
		}
	}
	return points, nil
}

func sparklineFor(u string, timeout time.Duration) (*sparkline, error) {
	safe := url.PathEscape(u)
	queryURL := fmt.Sprintf("%v/__wb/sparkline?url=%v&collection=web&output=json", BaseURL, safe)
	sl := &sparkline{}
	if _, err := simpleHTTPJSON(queryURL, sl, timeout); err != nil {
		return nil, err
	}
	sl.safeURL = safe
	sl.timeout = timeout
	return sl, nil
}

const maxRetries = 10

// simpleHTTPJSON deserializes response body content from get request url into
// objPtr.
func simpleHTTPJSON(u string, objPtr interface{}, timeout time.Duration, attempt ...int) (*http.Response, error) {
	if len(attempt) == 0 {
		attempt = []int{0}
	}
	attempt[0]++
	log.WithField("url", u).Debug("Downloading JSON data")
	resp, body, err := doRequest("", u, nil, timeout)
	if err != nil {
		u80 := strings.Replace(u, "https://", "http://", 1)
		log.Debug(u80)
		if resp, body, err = doRequest("", u80, nil, timeout); err != nil {
			log.Debugf("e2=%s", err)
			if attempt[0] > maxRetries {
				return resp, err
			}
			time.Sleep(5 * time.Second * time.Duration(attempt[0]))
			return simpleHTTPJSON(u, objPtr, timeout, attempt...)
		}
	}
	if err := json.Unmarshal(body, objPtr); err != nil {
		return resp, err
	}
	return resp, nil
}
