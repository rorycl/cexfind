package location

import (
	"fmt"
	"sort"
	"sync"
)

type StoreWithDistance struct {
	StoreID       int
	StoreName     string
	RegionName    string
	Latitude      float64
	Longitude     float64
	DistanceMiles float64
}

// StoreDistances finds the distances of the named stores from postcode
// and returns a slice of StoreWithDistance sorted by increasing
// distance
func StoreDistances(postcode string, storeNames []string) ([]StoreWithDistance, error) {

	location, err := getLocationFromPostcode(postcode)
	if err != nil {
		return nil, fmt.Errorf("could not resolve postcode: %w", err)
	}
	locationCoord := coord{Lat: location.Latitude, Lon: location.Longitude}

	// initialise the package global Stores
	var once sync.Once
	once.Do(func() {
		err = getStoreLocations()
	})
	if err != nil {
		return nil, fmt.Errorf("could not get store locations: %w", err)
	}

	foundStores := []StoreWithDistance{}
	for _, name := range storeNames {
		fs := StoreWithDistance{StoreName: name}
		thisStore, ok := Stores[name]
		if !ok {
			continue
		}
		fs.StoreID = thisStore.StoreID
		fs.RegionName = thisStore.RegionName
		fs.Latitude = thisStore.Latitude
		fs.Longitude = thisStore.Longitude

		// haversine distance in miles
		storeCoord := coord{Lat: fs.Latitude, Lon: fs.Longitude}
		fs.DistanceMiles, _ = haversineDistance(locationCoord, storeCoord) // mi, km

		foundStores = append(foundStores, fs)
	}

	sort.Slice(foundStores, func(i, j int) bool {
		return foundStores[i].DistanceMiles < foundStores[j].DistanceMiles
	})

	return foundStores, nil

}
