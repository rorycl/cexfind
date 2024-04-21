// Package cexfind searches for devices for sale at Cex/Webuy via the
// unofficial `webuy.io` query endpoint which responds in a json format.
//
// Queries are required to be made in the UK as the endpoint is
// protected by region-sensitive CDN.
//
// Example usage:
//
//    results, err := cex.Search(queries, strict)
//    if err != nil {
//    	log.Fatal(err)
//    }
//
//    for _, box := range results {
//    	fmt.Printf("%20s : %3d %s\n", box.Model, box.Price, box.Name)
//    }
package cexfind

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

// Box is a rationalised representation of a Cex/Webuy json entry, where
// each entry represents a "Box" or computer or other item for sale.
type Box struct {
	Model string
	Name  string
	ID    string
	Price int
}

// inQuery checks to see if each of the words in at least one of the
// supplied queries are in the Name of a Box. inQuery is used for
// determining if a particular Box should be returned from a "strict"
// Search
func (b *Box) inQuery(queries []string) bool {
	for _, q := range queries {
		matches := 0
		name := strings.ToLower(b.Name)
		words := strings.Split(strings.ToLower(q), " ")
		for _, w := range words {
			if strings.Contains(name, w) {
				matches++
			}
		}
		if matches == len(words) {
			return true
		}
	}
	return false
}

// IDUrl returns the full url path to the Cex/Webuy webpage showing Box
// in question
func (b Box) IDUrl() string {
	return urlDetail + b.ID
}

// boxes is a slice of Box
type boxes []Box

// sort boxes by a Box attribute
func (b boxes) sort(typer string) {
	sort.SliceStable(b, func(i, j int) bool {
		switch typer {
		case "ID":
			if b[i].ID < b[j].ID {
				return true
			}
		default: // sort by price
			if b[i].Price < b[j].Price {
				return true
			}
		}
		return false
	})
}

// boxMap is a map of boxes by model name, used for aggregating the
// results of several queries into a single map to avoid duplicate items
type boxMap map[string]boxes

// asBoxes returns an ordered slice of Box contained in the boxMap
func (b boxMap) asBoxes() []Box {
	var theseBoxes []Box
	keys := []string{}
	for k := range b {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		bSlice := b[k]
		bSlice.sort("Price")
		for _, b := range bSlice {
			theseBoxes = append(theseBoxes, b)
		}
	}
	return theseBoxes
}

// boxResults encapsulates the responses from a search query
type boxResults struct {
	boxmap boxMap
	err    error
}

// Search searches the Cex json endpoint at URL for the provided
// queries, returning a slice of Box or error. The strict flag ensures
// that the results contain terms from the search queries as the
// non-strict results include additional suggestions from the
// Cex/Webuy system.
func Search(queries []string, strict bool) ([]Box, error) {

	var allBoxes []Box
	allResults := boxMap{}

	// get chan results from the (potentially) multiple queries
	results := makeQueries(queries, strict)

	for br := range results {
		// exit on first error
		if br.err != nil {
			return allBoxes, br.err
		}

		// aggregate results and compact to remove duplicates
		for k, v := range br.boxmap {
			if _, ok := allResults[k]; !ok {
				allResults[k] = v
			} else {
				tmp := slices.Concat(allResults[k], v)
				tmp.sort("ID")
				allResults[k] = slices.Compact(tmp)
			}
		}
	}
	if len(allResults) == 0 {
		return allBoxes, errors.New("no results")
	}
	return allResults.asBoxes(), nil
}
