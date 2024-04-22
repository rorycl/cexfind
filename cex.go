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
	"errors"
	"slices"
	"sort"
	"strings"
)

// Box is a very simplified representation of a Cex/Webuy json entry,
// where each entry represents a "Box" or computer or other item for
// sale.
type Box struct {
	Model string
	Name  string
	ID    string
	Price int
}

// inQuery checks to see if each of the words in at least one of the
// supplied queries are in the Name of a Box. inQuery is used for
// determining if a particular Box should be returned from a "strict"
// search.
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

// IDUrl returns the full url path to the Cex/Webuy webpage showing the
// Box (the item of equipment) in question.
func (b Box) IDUrl() string {
	return urlDetail + b.ID
}

// boxes is a slice of Box
type boxes []Box

// sort boxes by box.model then box price ascending
func (b boxes) sort() {
	sort.SliceStable(b, func(i, j int) bool {
		bs, js := strings.ToLower(b[i].Model), strings.ToLower(b[j].Model)
		if bs != js {
			return bs < js
		}
		return b[i].Price < b[j].Price
	})
}

// Search searches the Cex json endpoint at URL for the provided
// queries, returning a slice of Box or error.
//
// The strict flag ensures that the results contain terms from the
// search queries as the non-strict results include additional
// suggestions from the Cex/Webuy system.
//
// Multiple queries are run concurrently and their results sorted by
// model, then by price ascending, and then aggregated to remove
// duplicates.
func Search(queries []string, strict bool) ([]Box, error) {
	var allBoxes boxes
	results := makeQueries(queries, strict)
	for br := range results {
		// exit on first error
		if br.err != nil {
			return allBoxes, br.err
		}
		allBoxes = append(allBoxes, br.box)
	}
	allBoxes.sort()
	allBoxes = slices.Compact(allBoxes)
	if len(allBoxes) == 0 {
		return allBoxes, errors.New("no results")
	}
	return allBoxes, nil
}
