// search for devices for sale at Cex
package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	//URL search url
	URL = "https://search.webuy.io/1/indexes/*/queries"
	// MODEL placeholder
	MODEL = "MODEL"
	// json body with placeholder; note that the availability online filter ensures only available kit is returned
	BODY = `{"requests":[{"indexName":"prod_cex_uk","params":"clickAnalytics=true&facetFilters=%5B%5B%22availability%3AIn%20Stock%20Online%22%5D%5D&facets=%5B%22*%22%5D&filters=boxVisibilityOnWeb%3D1%20AND%20boxSaleAllowed%3D1&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=17&maxValuesPerFacet=1000&page=0&query=MODEL&tagFilters=&userToken=71d182c769bd4dbc94081214a363c014"}]}`
)

// JsonResults encompasses the interesting fields in a Cex web search result
type JsonResults struct {
	Results []struct {
		Hits []struct {
			BoxName string `json:"boxName"`
			BoxID   string `json:"boxId"`
			// Available int `json:"collectionQuantity"` // returns 0 or greater
			Price int `json:"sellPrice"`
		} `json:"hits"`
		NbHits      int `json:"nbHits"`
		HitsPerPage int `json:"hitsPerPage"`
	} `json:"results"`
}

// Box is a rationalised JsonResults.Results.Hits entry, notionally
// representing a "Box" or computer or other item for sale
type Box struct {
	Model string
	Name  string
	ID    string
	Price int
}

// Boxes is a slice of Box
type Boxes []Box

// Sort Boxes by a Box attribute
func (b Boxes) Sort(typer string) {
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
	return
}

// BoxMap is a map of Boxes by model name
type BoxMap map[string]Boxes

// boxResults encapsulates the responses from a search query
type boxResults struct {
	boxmap BoxMap
	err    error
}

// headingExtract attempts to extract an h1 heading from a stream of
// bytes, typically needed if there is an html error page
func headingExtract(b []byte) string {
	reH1 := regexp.MustCompile(`<h1[^>]*>([^<]+)</h1>`)
	results := reH1.FindSubmatch(b)
	if len(results) < 2 {
		return ""
	}
	return string(results[1])
}

// postQuery posts the web query
func postQuery(queryBytes []byte) (JsonResults, error) {
	var r JsonResults
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(queryBytes))
	if err != nil {
		return r, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return r, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return r, err
	}

	/* save to a temporary json file for inspection */
	err = os.WriteFile("tmp.json", responseBytes, 0644)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		reason := headingExtract(responseBytes)
		if reason == "" {
			reason = "unknown or unmarshalling error"
		}
		return r, fmt.Errorf("error: %s", reason)
	}
	if len(r.Results) < 1 || len(r.Results[0].Hits) < 1 {
		return r, errors.New("no results")
	}
	return r, nil
}

// extractModelType tries to extract a meaningful model type from a
// boxname
func extractModelType(s string) string {
	parts := strings.Split(s, "/")
	if len(parts) > 1 {
		return strings.Title(strings.ToLower(parts[0]))
	}
	// grab first two words
	parts = strings.SplitN(s, " ", 3)
	if len(parts) == 1 {
		return strings.Title(strings.ToLower(parts[0]))
	}
	return strings.Title(strings.ToLower(strings.Join(parts[:2], " ")))
}

// makeQueries makes queries concurrently
func makeQueries(queries []string) chan boxResults {

	results := make(chan boxResults)

	go func() {
		defer close(results)

		for _, query := range queries {
			br := boxResults{}
			br.boxmap = BoxMap{}

			queryBody := strings.ReplaceAll(BODY, "MODEL", query)
			queryBytes := []byte(queryBody)

			response, err := postQuery(queryBytes)
			if err != nil {
				br.err = err
				results <- br
				return
			}

			for _, j := range response.Results[0].Hits {
				box := Box{}
				box.Model = extractModelType(j.BoxName)
				box.Name = j.BoxName
				box.ID = j.BoxID
				box.Price = j.Price

				if _, ok := br.boxmap[box.Model]; !ok {
					br.boxmap[box.Model] = []Box{}
				}
				br.boxmap[box.Model] = append(br.boxmap[box.Model], box)
			}
			results <- br
		}
	}()
	return results
}

// Search searches the Cex json endpoint at URL for the provided
// queries, returning a BoxMap or error
func Search(queries []string) (BoxMap, error) {

	allResults := BoxMap{}
	results := makeQueries(queries)

	for br := range results {
		// exit on first error
		if br.err != nil {
			return allResults, br.err
		}

		// aggregate results and compact to remove duplicates
		for k, v := range br.boxmap {
			if _, ok := allResults[k]; !ok {
				allResults[k] = v
			} else {
				tmp := slices.Concat(allResults[k], v)
				tmp.Sort("ID")
				allResults[k] = slices.Compact(tmp)
			}
		}
	}

	if len(allResults) == 0 {
		return allResults, errors.New("no results")
	}
	return allResults, nil
}
