package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var storeURL string = "https://wss2.cex.uk.webuy.io/v3/stores"

// storeLocations contains the interesting fields from the Cex store listings
type storeLocations struct {
	Response struct {
		Data struct {
			Stores []struct {
				StoreID    int     `json:"storeId"`
				StoreName  string  `json:"storeName"`
				RegionName string  `json:"regionName"`
				Latitude   float64 `json:"latitude"`
				Longitude  float64 `json:"longitude"`
				// PhoneNumber     any     `json:"phoneNumber"`
				ClosingTime string `json:"closingTime"`
			} `json:"stores"`
		} `json:"data"`
	} `json:"response"`
}

// Store is a store rationalised from storeLocations
type Store struct {
	StoreID    int
	StoreName  string
	RegionName string
	Latitude   float64
	Longitude  float64
}

// Stores is a map of Store by name
type Stores map[string]Store

func addAliases(s Stores) {
	simpleMap := map[string]string{
		"Tottenham Crt Rd": "London W1 TCR",
		"Rathbone Place":   "London W1 Rathbone",
	}
LOOP:
	for k, v := range s {
		for k2, v2 := range simpleMap {
			if strings.Contains(k, k2) {
				// make a new entry in the Stores map
				s[v2] = v
				continue LOOP
			}
		}
	}
}

func getLocations() (Stores, error) {

	stores := Stores{}

	var jsonStores storeLocations
	response, err := http.Get(storeURL)
	if err != nil {
		return stores, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return stores, fmt.Errorf("http response read error: %w", err)
	}

	err = json.Unmarshal(responseBytes, &jsonStores)
	if err != nil {
		return stores, fmt.Errorf("unmarshal error: %w", err)
	}

	for _, store := range jsonStores.Response.Data.Stores {
		// fmt.Printf("%3d %20s lat %5.8f long %5.8f\n", store.StoreID, store.StoreName, store.Latitude, store.Longitude)
		stores[store.StoreName] = Store{
			StoreID:    store.StoreID,
			StoreName:  store.StoreName,
			RegionName: store.RegionName,
			Latitude:   store.Latitude,
			Longitude:  store.Longitude,
		}
	}

	addAliases(stores)

	return stores, nil
}
