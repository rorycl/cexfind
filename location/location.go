package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// findLocationURL is the url for looking up location by UK postcode.
// Note that the full url is findLocationURL + ?q=
var findLocationURL string = "https://api.postcodes.io/postcodes"

var ErrLocationNotFound error = errors.New("location not found")

// Location represents the location derived from a web call
type location struct {
	Postcode  string
	District  string
	Latitude  float64
	Longitude float64
}

// jsonLocation is the raw data from a json geolocation call
type jsonLocation struct {
	Result []struct {
		Postcode  string  `json:"postcode"`
		Quality   int     `json:"quality"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		District  string  `json:"admin_district"`
	} `json:"result"`
}

// locationFinder holds information about postcodes, avoiding lookups
// for the same postcode
type locationFinder struct {
	client      *http.Client
	locationMap map[string]location
	sync.RWMutex
}

// newLocationFinder returns a new locationFinder. This should only be
// initialised once.
func newLocationFinder(client *http.Client) *locationFinder {
	l := locationFinder{
		client:      client,
		locationMap: map[string]location{},
	}
	return &l
}

func (lf *locationFinder) clean(postcode string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ToLower(postcode)), "  ", " ")
}

func (lf *locationFinder) has(postcode string) bool {
	lf.RLock()
	defer lf.RUnlock()
	if _, ok := lf.locationMap[lf.clean(postcode)]; ok {
		return true
	}
	return false
}

func (lf *locationFinder) get(postcode string) (location, bool) {
	lf.RLock()
	defer lf.RUnlock()
	l, ok := lf.locationMap[lf.clean(postcode)]
	return l, ok
}

func (lf *locationFinder) put(postcode string, l location) {
	lf.Lock()
	defer lf.Unlock()
	lf.locationMap[lf.clean(postcode)] = l
}

func (lf *locationFinder) length() int {
	lf.RLock()
	defer lf.RUnlock()
	return len(lf.locationMap)
}

type ErrorLocationHTTPGet struct {
	err error
}

func (elhg ErrorLocationHTTPGet) Error() string {
	return elhg.err.Error()
}

// getLocationFromPostcode tries to extract the location data from the
// locationMap. If that fails it tries to extract the necessary
// information from a web service.
func (lf *locationFinder) getLocationFromPostcode(postcode string) (*location, error) {

	if postcode == "" {
		return nil, errors.New("no postcode provided")
	}

	// check postcode in cache
	if l, ok := lf.get(postcode); ok {
		return &l, nil
	}

	var jloc jsonLocation

	response, err := lf.client.Get(findLocationURL + "?q=" + url.QueryEscape(postcode))
	if err != nil {
		return nil, ErrorLocationHTTPGet{err}
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("http response read error: %w", err)
	}

	err = json.Unmarshal(responseBytes, &jloc)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	if len(jloc.Result) < 1 {
		return nil, ErrLocationNotFound
	}

	result := jloc.Result[0]

	l := location{
		Postcode:  result.Postcode,
		District:  result.District,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}
	lf.put(postcode, l)
	return &l, nil
}
