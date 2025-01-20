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
//	results, err := cex.Search(queries, strict)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	latestModel := ""
//	for _, box := range results {
//		if box.Model != latestModel {
//			fmt.Printf("\n%s\n", box.Model)
//			latestModel = box.Model
//		}
//		fmt.Printf(
//			"   £%3d %s\n   %s\n",
//			box.Price,
//			box.Name,
//			box.IDUrl(),
//		)
//	}
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
// roughly length, truncating with "..." where necessary
func (b *Box) StoresString(length int) string {
	storeString := ""
	if len(b.Stores) > 0 {
		storeString += fmt.Sprintf("%s", b.Stores[0])
	}
	if len(b.Stores) < 2 {
		return storeString
	}
	for _, s := range b.Stores[1:] {
		storeString += fmt.Sprintf(", %s", s)
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

// Search searches the Cex json endpoint at URL for the provided
// queries, returning a slice of Box or error.
//
// The strict flag ensures that the results contain terms from the
// search queries as the non-strict results include additional
// suggestions from the Cex/Webuy system.
//
// The postcode, if provided, allows distances to be calculated from
// each store.
//
// Multiple queries are run concurrently and their results sorted by
// model, then by price ascending. Duplicate results are removed at
// aggregation.
func Search(queries []string, strict bool, postcode string) ([]Box, error) {
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

		br.box.Stores, err = location.StoreDistances(postcode, br.box.storeNames)
		if err != nil {
			err = fmt.Errorf("Postcode error: %w", err)
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
