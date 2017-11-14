package webex

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
)

// Scrape scrapes a web page for the regular expression area,
// then the area for the regular expression spot.
func Scrape(link, area, spot string) ([]byte, error) {
	areaRegex, err := regexp.Compile(area)
	if err != nil {
		return nil, err
	}
	spotRegex, err := regexp.Compile(spot)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	areas := areaRegex.FindAll(b, 1)
	if len(areas) == 0 {
		return nil, errors.New("no area match")
	}
	spots := spotRegex.FindAll(areas[0], -1)
	if len(spots) == 0 {
		return nil, errors.New("no spot match")
	}
	return bytes.Join(spots, nil), nil
}

// ScrapeString scrapes a web page for the regular expression area,
// then the area for the regular expression spot into a string.
func ScrapeString(link, area, spot string) (string, error) {
	b, err := Scrape(link, area, spot)
	return string(b), err
}

// ScrapeInt scrapes a web page for the regular expression area,
// then the area for the regular expression spot into an int.
func ScrapeInt(link, area, spot string) (int, error) {
	b, err := Scrape(link, area, spot)
	return int(binary.BigEndian.Uint64(b)), err
}

// ScrapeFloat64 scrapes a web page for the regular expression area,
// then the area for the regular expression spot into a float64.
func ScrapeFloat64(link, area, spot string) (float64, error) {
	b, err := Scrape(link, area, spot)
	return math.Float64frombits(binary.BigEndian.Uint64(b)), err
}
