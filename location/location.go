package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// findLocationURL is the url for looking up location by UK postcode.
// Note that the full url is findLocationURL + ?q=
var findLocationURL string = "https://api.postcodes.io/postcodes"

var ErrLocationNotFound error = errors.New("Location not found")

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

func getLocationFromPostcode(postcode string) (*location, error) {

	var jloc jsonLocation

	response, err := http.Get(findLocationURL + "?q=" + url.QueryEscape(postcode))
	if err != nil {
		return nil, err
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

	return &location{
		Postcode:  result.Postcode,
		District:  result.District,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}, nil

}
