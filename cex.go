// Package cexfind searches for devices for sale at Cex/Webuy via the
// unofficial "webuy.io" query endpoint which responds in a json format.
//
// Queries are required to be made in the UK as the endpoint is
// protected by region-sensitive CDN.
//
// Multiple concurrent queries are supported, with an optional "strict"
// flag to constrain results to the query terms. The results are a union
// of the results of each query, ordered by model name and then the
// price of each item.
//
// Example usage:
//
//		queries := []string{"Lenovo T14s"}
//		strict := false
//		postcode := "S10 1LT" // royal armouries museum, leeds
//		var kit *CexFind
//		if postcode != "" {
//			kit, _ := cex.NewCexFind(WithStoreDistanceInitiliase())
//		} else {
//			kit, _ := cex.NewCexFind()
//	    }
//		results, err := cex.Search(queries, strict, postcode)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		for _, box := range results {
//			fmt.Printf("\n%s\n", box.Model)
//			fmt.Printf(
//			"   £%3d %s\n   %s\n   %s\n",
//				box.Price,
//				box.Name,
//				box.IDUrl(),
//				box.StoresString(100), // up to 100 chars of store info
//			)
//		}
package cexfind

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/rorycl/cexfind/location"
	"github.com/shopspring/decimal"
)

// Box is a very simplified representation of a Cex/Webuy json entry,
// where each entry represents a "Box" or computer or other item for
// sale.
type Box struct {
	Model         string
	Name          string
	Category      string
	ID            string
	Price         decimal.Decimal
	PriceCash     decimal.Decimal
	PriceExchange decimal.Decimal
	storeNames    []string
	Stores        []location.StoreWithDistance
}

// inQuery checks to see if each of the words in at least one of the
// supplied queries are in the Name of a Box. inQuery is used for
// determining if a particular Box should be returned from a "strict"
// search.
func (b *Box) inQuery(queries []string) bool {
	for _, q := range queries {
		matches := 0
		name := strings.ToLower(b.Name + " " + b.Model)
		words := strings.Split(strings.ToLower(q), " ")
		for _, w := range words {
			if strings.Contains(name, w) {
				matches++
				continue
			}
		}
		if matches == len(words) {
			return true
		}
	}
	return false
}

// IDUrl returns the full url path to the Cex/Webuy webpage showing the
// Box (the item of equipment) in question.
func (b Box) IDUrl() string {
	return urlDetail + b.ID
}

// reverseID is useful for sorting because the grade of the box is the
// right-most character. The grade cannot be conveniently extracted
// otherwise. For the same price, a higher grade (eg B) is prefereable
// over a lower grade (eg C).
func (b *Box) reverseID() string {
	r := []rune(b.ID)
	for i, j := 0, len(b.ID)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// StoresString returns the stores as a comma delimited string to
// roughly length, truncating with "…" where necessary. Giving
// StoreString an argument length of -1 means there is no limit on the
// length of the returned string.
func (b *Box) StoresString(length int) string {

	storeString := ""
	if len(b.Stores) == 0 {
		return storeString
	}

	storeString += b.Stores[0].String()
	for _, s := range b.Stores[1:] {
		storeString += fmt.Sprintf(", %s", s)
	}

	if len(b.Stores) < 2 || length == -1 {
		return storeString
	}
	if len(storeString) > length {
		storeString = storeString[:length]
		storeString = strings.TrimSuffix(storeString, " ")
		storeString = strings.TrimSuffix(storeString, ",")
		storeString = storeString + "…"
	}
	return storeString
}

// boxes is a slice of Box
type boxes []Box

// sort sorts boxes by box.Model then Price ascending then ID.
func (b *boxes) sort() {
	slices.SortFunc(*b, func(i, j Box) int {
		var c int
		c = cmp.Compare(i.Model, j.Model)
		if c != 0 {
			return c
		}
		c = i.Price.Compare(j.Price)
		if c != 0 {
			return c
		}
		// the most right char is the box condition (A, B or C)
		return cmp.Compare(i.reverseID(), j.reverseID())
	})
}

// CexFind provides the means for searching Cex's API with store
// location data and a configurable http Client.
type CexFind struct {
	client         *http.Client
	storeDistances *location.StoreDistances
}

// Option is an initialiser opt.
type Option func(*CexFind) error

// NewCexFind makes a new Cex instance taking zero or more Option. To use the
// `WithProxy` for example the function should be called along the following lines:
//
//	c, err := NewCexFind(WithProxy("socks5://127.0.0.1:8081"))
//
// If store distances from a postcode are needed, use the WithStoreDistanceInitiliase
// option. For store distances with periodic updates (for long running programmes, such
// as the web interface perhaps), use the WithStoreDistanceInitializeAndUpdates option.
// Note that only one of WithStoreDistanceInitiliase and
// WithStoreDistanceInitializeAndUpdates should be used.
//
// An example of using both a socks5 proxy and WithStoreDistanceInitiliase:
//
//	c, err := NewCexFind(
//		WithProxy("socks5://127.0.0.1:8081"),
//		WithStoreDistanceInitialize(),
//	)
func NewCexFind(opts ...Option) (*CexFind, error) {
	client, _ := newHTTPClient("")

	// Setup the stores and location finder.
	nsd, err := location.NewStoreDistances(client)
	if err != nil {
		return nil, fmt.Errorf("NewStoreDistances new error: %w", err)
	}

	cex := &CexFind{
		client:         client,
		storeDistances: nsd,
	}

	// Add options, if any.
	for _, opt := range opts {
		if err := opt(cex); err != nil {
			return nil, err
		}
	}
	return cex, nil
}

// proxySchemeOK checks if the proxy url includes the expected `scheme` component.
func proxySchemeOK(u *url.URL) error {
	expected := []string{"socks5", "http", "https"}
	for _, mtch := range expected {
		if strings.Contains(u.Scheme, mtch) {
			return nil
		}
	}
	return fmt.Errorf("proxy schema %q in %s not recognised, expected one of %v", u.Scheme, u.String(), expected)
}

// newHTTPClient is a convenience func for initialising an http client from scratch at
// initialisation or update ("with" options), since updating a transport via
// <obj>.(*http.Transport) appears unreliable.
func newHTTPClient(proxy string) (*http.Client, error) {

	maxConnsPerHost := 8
	timeout := time.Second * 15

	if proxy == "" {
		return &http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost: maxConnsPerHost,
			},
			Timeout: timeout,
		}, nil
	}

	// Parse and check the provided proxy url.
	u, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("proxy %q could not be parsed: %w", proxy, err)
	}
	if err := proxySchemeOK(u); err != nil {
		return nil, err
	}

	// Attach a new client.
	return &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: maxConnsPerHost,
			Proxy:           http.ProxyURL(u),
		},
		Timeout: timeout,
	}, nil
}

// WithProxy configures the HTTP client to use a proxy (e.g. socks5://...)
func WithProxy(proxyStr string) Option {
	return func(c *CexFind) error {
		var err error
		c.client, err = newHTTPClient(proxyStr)
		return err
	}
}

// WithStoreDistanceInitialize initialises the store distancer.
func WithStoreDistanceInitiliase() Option {
	return func(c *CexFind) error {
		if c.storeDistances.IsOperational() {
			return errors.New("store distances already initialised")
		}
		err := c.storeDistances.Initialise()
		if err != nil {
			return fmt.Errorf("store distances initialisation error: %w", err)
		}
		return nil
	}
}

// WithStoreDistanceInitializeAndUpdates initialises and starts running the periodic
// store distance updater.
func WithStoreDistanceInitializeAndUpdates() Option {
	return func(c *CexFind) error {
		if c.storeDistances.IsOperational() {
			return errors.New("store distances already initialised")
		}
		err := c.storeDistances.Initialise()
		if err != nil {
			return fmt.Errorf("store distances initialisation error: %w", err)
		}
		c.storeDistances.StartPeriodicReinitialisation()
		return nil
	}
}

// LocationDistancesOK indicates if the storeDistances.store has been
// initialised and distances can be calculated. If the func returns
// false then store distances won't be calculated, a fact that client
// apps should probably report.
func (c *CexFind) LocationDistancesOK() bool {
	return c.storeDistances.IsOperational()
}

// Search searches the Cex json endpoint at URL for the provided
// queries, returning a slice of Box or error.
//
// The strict flag ensures that the results contain terms from the
// search queries as the non-strict results include additional
// suggestions from the Cex/Webuy system.
//
// The postcode, if provided, allows distances to be calculated from
// each store if the store data has already been retrieved (store data
// is retrieved asynchronously).
//
// Multiple queries are run concurrently and their results sorted by
// model, then by price ascending. Duplicate results are removed at
// aggregation.
func (cex *CexFind) Search(queries []string, strict bool, postcode string) ([]Box, error) {
	var allBoxes boxes
	var idMap = make(map[string]struct{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	results := makeQueries(ctx, cex.client, queries, strict)

	for br := range results {
		if br.err != nil {
			if err != nil {
				err = fmt.Errorf("\"%s\": %w\n%w", br.query, br.err, err)
			} else {
				err = fmt.Errorf("\"%s\": %w", br.query, br.err)
			}
			continue
		}
		if _, ok := idMap[br.box.ID]; ok { // don't add duplicates
			continue
		}

		// Store information is cached, as is any postcode with its
		// location data. cached data only requires distances to be
		// calculated. If stores are offline distance calcs are skipped,
		// but stores "with distances" are still returned.
		br.box.Stores, err = cex.storeDistances.Distances(postcode, br.box.storeNames)
		if err != nil {
			err = fmt.Errorf("postcode error: %w", err)
			return nil, err
		}

		allBoxes = append(allBoxes, br.box)
		idMap[br.box.ID] = struct{}{}
	}
	allBoxes.sort()
	if len(allBoxes) == 0 {
		if err != nil {
			err = fmt.Errorf("%w", err)
		} else {
			err = errors.New("no results")
		}
		return allBoxes, err
	}
	return allBoxes, err
}
