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
// cache is updated by default once a day.
type stores struct {
	URL         string
	client      *http.Client
	storeMap    map[string]store
	initialised bool
	update      *time.Ticker
	sync.RWMutex
}

var tickerOKDuration time.Duration = time.Minute * 60 * 24
var tickerProblemDuration time.Duration = time.Minute * 10

// newStores initialises a concurrent safe stores struct. Normally this is initialised
// with an http client initialised and shared from a caller.
func newStores(client *http.Client) *stores {
	if client == nil {
		client = http.DefaultClient
	}
	s := stores{
		URL:      storeURL, // default
		client:   client,
		storeMap: map[string]store{},
		update:   time.NewTicker(tickerOKDuration),
	}
	return &s
}

// initialise initialises the store database, returning an error if this fails.
func (s *stores) initialise() error {
	err := s.getStoreLocations()
	if err != nil {
		return fmt.Errorf("store database initialisation error %s", err)
	} else {
		s.Lock()
		s.initialised = true
		s.Unlock()
	}
	return nil
}

// periodicallyReinitialise periodicically updates for long running clients, such as the
// web client.
func (s *stores) periodicallyReinitialise() {
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
}

// get the store information for a named store.
func (s *stores) get(name string) (*store, bool) {
	if !s.isInitialised() {
		return nil, false
	}

	s.RLock()
	defer s.RUnlock()
	st, ok := s.storeMap[name]
	return &st, ok
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

func (s *stores) list() []*store {
	var all []*store
	s.Lock()
	defer s.Unlock()
	for _, v := range s.storeMap {
		all = append(all, &v)
	}
	return all
}

// getStoreLocations gets the store locations from the storeURL and
// processes them into the stores map by the store name.
func (s *stores) getStoreLocations() error {

	var jsonStores storeLocations

	response, err := s.client.Get(s.URL)
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
		// fmt.Printf("%3d %20s lat %5.8f long %5.8f\n", jStore.StoreID, jStore.StoreName, jStore.Latitude, jStore.Longitude)
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
