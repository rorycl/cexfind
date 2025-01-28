package location

import (
	"fmt"
	"sort"
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

// StoreDistances is a struct containing a map of stores by name
// allowing calculations for distances from a postcode.
type StoreDistances struct {
	stores         *stores
	locationFinder *locationFinder
}

// NewStoreDistances initialises a StoresDistance instance, and passes a
// testing flag to the stores initaliser. In production,
// initialiseStores should be true
func NewStoreDistances(initialiseStores bool) *StoreDistances {
	s := StoreDistances{
		stores:         newStores(initialiseStores),
		locationFinder: newLocationFinder(),
	}
	return &s
}

// Distances finds the distances of the named stores from postcode and
// returns a slice of StoreWithDistance sorted by increasing distance
func (sd *StoreDistances) Distances(postcode string, storeNames []string) ([]StoreWithDistance, error) {

	foundStores := []StoreWithDistance{}

	// return sparse stores if no postcode is provided or stores haven't
	// been initialised.
	if postcode == "" || !sd.IsOperational() {
		foundStores := []StoreWithDistance{}
		for _, name := range storeNames {
			swd := StoreWithDistance{StoreName: name}
			foundStores = append(foundStores, swd)
		}
		storeSorter(foundStores)
		return foundStores, nil
	}

	location, err := sd.locationFinder.getLocationFromPostcode(postcode)
	if err != nil {
		return nil, fmt.Errorf("could not resolve postcode: %w", err)
	}
	locationCoord := coord{Lat: location.Latitude, Lon: location.Longitude}

	// hope the API uses consistent store names
	for _, name := range storeNames {
		fs := StoreWithDistance{StoreName: name}
		thisStore, ok := sd.stores.get(name)
		if ok { // allow sparse stores
			fs.StoreID = thisStore.StoreID
			fs.RegionName = thisStore.RegionName
			fs.Latitude = thisStore.Latitude
			fs.Longitude = thisStore.Longitude

			// haversine distance in miles
			storeCoord := coord{Lat: fs.Latitude, Lon: fs.Longitude}
			fs.DistanceMiles, _ = haversineDistance(locationCoord, storeCoord) // mi, km
		}
		foundStores = append(foundStores, fs)
	}
	storeSorter(foundStores)
	return foundStores, nil
}

// IsOperational determines if the stores have been initalised and
// therefore if distances are possible to be calculated
func (sd *StoreDistances) IsOperational() bool {
	return sd.stores.isInitialised()
}
