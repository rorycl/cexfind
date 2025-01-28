package location

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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
type store struct {
	StoreID    int
	StoreName  string
	RegionName string
	Latitude   float64
	Longitude  float64
}

// stores is a collection of store safe for concurrent access. The Store
// cache is updated once a day.
type stores struct {
	storeMap map[string]store
	sync.RWMutex
	initialised bool
	update      *time.Ticker
}

var tickerOKDuration time.Duration = time.Minute * 60 * 24
var tickerProblemDuration time.Duration = time.Minute * 10

// newStores initialises a concurrent safe stores struct. The stores are
// only initialised if true, which is the default in production.
func newStores(initialiseStores bool) *stores {
	s := stores{
		storeMap: map[string]store{},
		update:   time.NewTicker(tickerOKDuration),
	}
	if initialiseStores {
		err := s.getStoreLocations()
		if err != nil {
			log.Printf("store update error %s", err)
			s.update.Reset(tickerProblemDuration)
		} else {
			s.initialised = true
		}
	}
	go func() {
		for range s.update.C {
			err := s.getStoreLocations()
			if err != nil {
				s.Lock()
				s.update.Reset(tickerProblemDuration)
				s.Unlock()
				log.Printf("store update error %s", err)
			} else {
				s.Lock()
				log.Println("store updated")
				s.initialised = true
				s.update.Reset(tickerOKDuration)
				s.Unlock()
			}
		}
	}()
	return &s
}

func (s *stores) get(name string) (store, bool) {
	s.RLock()
	defer s.RUnlock()
	st, ok := s.storeMap[name]
	return st, ok
}

func (s *stores) isInitialised() bool {
	s.RLock()
	defer s.RUnlock()
	return s.initialised
}

func (s *stores) length() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.storeMap)
}

func (s *stores) addAliases() {
	simpleMap := map[string]string{
		"Tottenham Crt Rd": "London W1 TCR",
		"Rathbone Place":   "London W1 Rathbone",
	}
	s.Lock()
	defer s.Unlock()
LOOP:
	for k, v := range s.storeMap {
		for k2, v2 := range simpleMap {
			if strings.Contains(k, k2) {
				// make a new entry in the stores map
				s.storeMap[v2] = v
				continue LOOP
			}
		}
	}
}

// getStoreLocations gets the store locations from the storeURL and
// processes them into the stores map by the store name.
func (s *stores) getStoreLocations() error {

	var jsonStores storeLocations
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	response, err := client.Get(storeURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("http response read error: %w", err)
	}

	err = json.Unmarshal(responseBytes, &jsonStores)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	s.Lock()
	for _, jStore := range jsonStores.Response.Data.Stores {
		// fmt.Printf("%3d %20s lat %5.8f long %5.8f\n", store.StoreID, store.StoreName, store.Latitude, store.Longitude)
		s.storeMap[jStore.StoreName] = store{
			StoreID:    jStore.StoreID,
			StoreName:  jStore.StoreName,
			RegionName: jStore.RegionName,
			Latitude:   jStore.Latitude,
			Longitude:  jStore.Longitude,
		}
	}
	s.Unlock()
	s.addAliases()
	return nil
}
