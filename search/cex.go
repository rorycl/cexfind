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

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//URL search url
	URL = "https://search.webuy.io/1/indexes/*/queries"
	// MODEL placeholder
	MODEL = "MODEL"
	// json body with placeholder; note that the availability online filter ensures only available kit is returned
	BODY = `{"requests":[{"indexName":"prod_cex_uk","params":"clickAnalytics=true&facetFilters=%5B%5B%22availability%3AIn%20Stock%20Online%22%5D%5D&facets=%5B%22*%22%5D&filters=boxVisibilityOnWeb%3D1%20AND%20boxSaleAllowed%3D1&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=17&maxValuesPerFacet=1000&page=0&query=MODEL&tagFilters=&userToken=71d182c769bd4dbc94081214a363c014"}]}`
	// save web output to temp file if DEBUG true
	DEBUG = false
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

// BoxMap is a map of boxes by model name
type BoxMap map[string]boxes

// BoxMapIter is a BoxMap key/Box pair
type boxMapIter struct {
	Key string
	Box Box
}

// Iter iterates over a BoxMap returning a key, Box pair over the whole
// BoxMap collection
func (b BoxMap) Iter() <-chan boxMapIter {
	bmi := make(chan boxMapIter)

	go func() {
		defer close(bmi)
		keys := []string{}
		for k := range b {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for _, k := range keys {
			v := b[k]
			v.sort("Price")
			for _, iv := range v {
				bmi <- boxMapIter{k, iv}
			}
		}
	}()
	return bmi
}

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
	if DEBUG {
		err = os.WriteFile("tmp.json", responseBytes, 0644)
		if err != nil {
			panic(err)
		}
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
// boxname. Since models are not well normalised further cleaning work
// is likely to be needed in future.
func extractModelType(s string) string {

	// words to remove from box model descriptions
	stringsToRemove := strings.Split("Thinkpad AddMoreHere", " ")

	// set the titling type to English
	titleCase := cases.Title(language.English)

	// clean some Title strings to remove string r + space
	cleaner := func(o, r string) string {
		i := strings.Index(strings.ToLower(o), strings.ToLower(r))
		if i == -1 {
			return o
		}
		return strings.ReplaceAll(o[:i]+o[i+len(r):], "  ", " ")
	}
	// remove known
	cleaners := func(o string) string {
		for _, w := range stringsToRemove {
			o = cleaner(o, w)
		}
		return o
	}

	// some titles have a "/" character summarizing the model
	parts := strings.Split(s, "/")
	if len(parts) > 1 {
		return cleaners(titleCase.String(strings.ToLower(parts[0])))
	}
	// grab first two words
	parts = strings.SplitN(s, " ", 3)
	if len(parts) == 1 {
		return titleCase.String(strings.ToLower(parts[0]))
	}
	return cleaners(titleCase.String(strings.ToLower(strings.Join(parts[:2], " "))))
}

// makeQueries makes queries concurrently; strict true requires that the
// return results contain all terms in at least one query
func makeQueries(queries []string, strict bool) chan boxResults {

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

				// in strict mode, don't add box if it doesn't match any query
				if strict && !box.inQuery(queries) {
					continue
				}

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
// queries, returning a BoxMap or error. The strict flag ensures that
// the results contain terms from the search queries
func Search(queries []string, strict bool) (BoxMap, error) {

	allResults := BoxMap{}
	results := makeQueries(queries, strict)

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
				tmp.sort("ID")
				allResults[k] = slices.Compact(tmp)
			}
		}
	}

	if len(allResults) == 0 {
		return allResults, errors.New("no results")
	}
	return allResults, nil
}
