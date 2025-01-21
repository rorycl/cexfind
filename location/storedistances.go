package location

import (
	"fmt"
	"sort"
	"sync"
)

// StoreWithDistance represents a store with a distance DistanceMiles
// from the provided postcode.
type StoreWithDistance struct {
	StoreID       int
	StoreName     string
	RegionName    string
	Latitude      float64
	Longitude     float64
	DistanceMiles float64
}

func (s StoreWithDistance) String() string {
	tplEmpty := "%s"
	tplShort := "%s (%.1fmi)"
	tplLong := "%s (%.fmi)"
	if s.StoreID == 0 {
		return fmt.Sprintf(tplEmpty, s.StoreName)
	}
	if s.DistanceMiles <= 10 {
		return fmt.Sprintf(tplShort, s.StoreName, s.DistanceMiles)
	}
	return fmt.Sprintf(tplLong, s.StoreName, s.DistanceMiles)
}

// storeSorter sorts a slice of StoreWithDistance, pulled out for
// testing
func storeSorter(fss []StoreWithDistance) {
	sort.Slice(fss, func(i, j int) bool {
		if fss[i].DistanceMiles == fss[j].DistanceMiles {
			return fss[i].StoreName < fss[j].StoreName
		}
		return fss[i].DistanceMiles < fss[j].DistanceMiles
	})
}

// StoreDistances finds the distances of the named stores from postcode
// and returns a slice of StoreWithDistance sorted by increasing
// distance
func StoreDistances(postcode string, storeNames []string) ([]StoreWithDistance, error) {

	foundStores := []StoreWithDistance{}

	// return sparse stores if no postcode is provided
	if postcode == "" {
		foundStores := []StoreWithDistance{}
		for _, name := range storeNames {
			swd := StoreWithDistance{StoreName: name}
			foundStores = append(foundStores, swd)
		}
		storeSorter(foundStores)
		return foundStores, nil
	}

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
	storeSorter(foundStores)

	return foundStores, nil

}
