package location

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var originalURL string = findLocationURL

func TestGetLocationLocal(t *testing.T) {

	testdata, err := os.ReadFile("testdata/findlocation.json")
	if err != nil {
		t.Fatal(err)
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(testdata))
	}))

	// repoint url
	findLocationURL = svr.URL

	postcode := "NW1 6LG"
	lFinder := newLocationFinder(&http.Client{Timeout: time.Second * 2})
	location, err := lFinder.getLocationFromPostcode(postcode)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := location.District, "Westminster"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
	if got, want := location.Latitude, 51.523969; got != want {
		t.Errorf("latitude got %f want %f", got, want)
	}
	if got, want := location.Longitude, -0.166312; got != want {
		t.Errorf("longitude got %f want %f", got, want)
	}

	if !lFinder.has(postcode) {
		t.Error("postcode could not be found in cache")
	}
	if got, want := lFinder.length(), 1; got != want {
		t.Errorf("cache map len got %d want %d", got, want)
	}
}

// TestGetLocationReal uses a live callout.
func TestGetLocationReal(t *testing.T) {

	findLocationURL = originalURL

	client := &http.Client{
		Transport: &http.Transport{
			// Proxy:           http.ProxyURL(someProxyURL),
		},
		Timeout: time.Millisecond * 750,
	}

	// Marlborough Mound 51.4166° N, 1.7371° W (numbers below are a bit
	// off)
	postcode := "SN8 1PA"
	lFinder := newLocationFinder(client)
	location, err := lFinder.getLocationFromPostcode(postcode)
	if _, ok := errors.AsType[ErrorLocationHTTPGet](err); ok {
		t.Fatalf("http error\n** This test requires internet access **\n%s", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	// location.Location{Postcode:"SN8 1PA", District:"Wiltshire", Latitude:51.417369, Longitude:-1.735758}

	if got, want := location.District, "Wiltshire"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
	if got, want := location.Latitude, 51.417367; got != want {
		t.Errorf("latitude got %f want %f", got, want)
	}
	if got, want := location.Longitude, -1.735748; got != want {
		t.Errorf("longitude got %f want %f", got, want)
	}
	if !lFinder.has(postcode) {
		t.Error("postcode could not be found in cache")
	}
	if got, want := lFinder.length(), 1; got != want {
		t.Errorf("cache map len got %d want %d", got, want)
	}
}
