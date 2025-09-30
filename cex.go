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
//	 postcode := "S10 1LT" // royal armouries museum, leeds
//		kit := cex.NewCex()
//		cex.Search(queries, strict, postcode)
//		results, err := kit.Search(queries, strict)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		latestModel := ""
//		for _, box := range results {
//			if box.Model != latestModel {
//				fmt.Printf("\n%s\n", box.Model)
//				latestModel = box.Model
//			}
//			fmt.Printf(
//				"   £%3d %s\n   %s\n   %s\n",
//				box.Price,
//				box.Name,
//				box.IDUrl(),
//				box.StoresString(100), // up to 100 chars of store info
//			)
//		}
package cexfind

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

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
// location data.
type CexFind struct {
	storeDistances *location.StoreDistances
}

// NewCexFind makes a new Cex instance. This should only be initalised
// once due to caching in the location submodules.
func NewCexFind() *CexFind {
	return &CexFind{
		storeDistances: location.NewStoreDistances(true),
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

	var err error

	results := makeQueries(queries, strict)
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
