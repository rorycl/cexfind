// query provides the web query functions, json marshalling and some
// content cleaning/management functions

package cexfind

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// URL is the Cex/Webuy search endpoint
	URL = "https://search.webuy.io/1/indexes/*/queries"
	// json body with placeholder MODEL; note that the availability online filter ensures only available kit is returned
	jsonBody = strings.ReplaceAll(`{"requests": [
    {
      "indexName": "prod_cex_uk",
      "params": "clickAnalytics=true
		&facetFilters=%5B%5B%22availability%3AIn%20Stock%20Online%22%5D%5D
		&facets=%5B%22*%22%5D
		&filters=boxVisibilityOnWeb%3D1%20
		&highlightPostTag=__%2Fais-highlight__
		&highlightPreTag=__ais-highlight__
		&hitsPerPage=50
		&maxValuesPerFacet=1000
		&page=0
		&query=MODEL
		&tagFilters=
		&userToken=71d182c769bd4dbc94081214a363c014"
    }]}`, "\n		", "")

	// urlDetail is the Cex/Webuy base url for individual items
	urlDetail = "https://uk.webuy.com/product-detail?id="
	// save web output to temp file if DEBUG true
	debug = false
	// no results sentinel error
	NoResultsFoundError error = errors.New("no results found")
)

// jsonResults encompasses the interesting fields in a Cex web search result
type jsonResults struct {
	Results []struct {
		Hits []struct {
			BoxName       string          `json:"boxName"`
			BoxID         string          `json:"boxId"`
			Category      string          `json:"categoryFriendlyName"`
			Price         decimal.Decimal `json:"sellPrice"`
			PriceCash     decimal.Decimal `json:"cashPriceCalculated"`     // offer price for this kit in cash
			PriceExchange decimal.Decimal `json:"exchangePriceCalculated"` // offer price for exchange
			Stores        []string        `json:"stores"`
		} `json:"hits"`
		// NbHits      int `json:"nbHits",omitempty`
		// HitsPerPage int `json:"hitsPerPage",omitempty`
	} `json:"results"`
}

// boxResults encapsulates the responses from a search query
type boxResults struct {
	query string
	box   Box
	err   error
}

// makeQueries makes queries concurrently; strict true requires that the
// return results contain all terms in at least one query
func makeQueries(queries []string, strict bool) chan boxResults {
	results := make(chan boxResults)
	var wg sync.WaitGroup
	for _, query := range queries {
		wg.Add(1)
		go func() {
			defer wg.Done()
			br := boxResults{query: query}
			queryBody := strings.ReplaceAll(jsonBody, "MODEL", url.QueryEscape(query))
			queryBytes := []byte(queryBody)
			response, err := postQuery(queryBytes)
			if err != nil {
				br.err = err
				results <- br
				return
			}

			for _, j := range response.Results[0].Hits {
				br.box = Box{}
				br.box.Model = extractModelType(j.BoxName)
				br.box.Name = j.BoxName
				br.box.Category = j.Category
				br.box.ID = j.BoxID
				br.box.Price = j.Price
				br.box.PriceCash = j.PriceCash
				br.box.PriceExchange = j.PriceExchange
				br.box.storeNames = storeSimplifier(j.Stores)
				// in strict mode, don't add box if it doesn't match any query
				if strict && !br.box.inQuery(queries) {
					continue
				}
				results <- br
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

// postQuery posts the web query
func postQuery(queryBytes []byte) (jsonResults, error) {
	var r jsonResults
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(queryBytes))
	if err != nil {
		return r, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return r, fmt.Errorf("http call error: %w", err)
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return r, fmt.Errorf("http response read error: %w", err)
	}

	/* save to a temporary json file for inspection */
	if debug {
		err = os.WriteFile("tmp.json", responseBytes, 0600)
		if err != nil {
			panic(err)
		}
	}

	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		// html page might have been returned; try and extract heading
		reason := errorExtract(responseBytes)
		if reason != "" {
			reason = "Cex API unknown retrieval or unmarshalling error"
			return r, errors.New(reason)
		}
		var ju *json.UnmarshalTypeError
		if errors.As(err, &ju) {
			return r, fmt.Errorf("json unmarshalling error: %w", err)
		}
		return r, fmt.Errorf("json reading error %w", err)
	}
	if len(r.Results) < 1 || len(r.Results[0].Hits) < 1 {
		return r, NoResultsFoundError
	}
	return r, nil
}

// errorExtract attempts to extract meaningful error messages from html
// error pages. It looks for an h1 heading from a stream of bytes,
// typically needed if there is an html error page, alternatively a
// CloudFlare blocking report.
func errorExtract(b []byte) string {
	reH1 := regexp.MustCompile(`<h1[^>]*>([^<]+)</h1>`)
	results := reH1.FindSubmatch(b)
	if len(results) >= 2 {
		return string(results[1])
	}
	if bytes.Contains(b, []byte("This website is using a security service to protect itself from online attacks")) {
		return "CloudFlare has blocked this service."
	}
	return ""
}

// storeSimplifier replaces some long store names with shorter ones
func storeSimplifier(ss []string) []string {
	output := []string{}
	simpleMap := map[string]string{
		"Tottenham Crt Rd": "London W1 TCR",
		"Rathbone Place":   "London W1 Rathbone",
	}
LOOP:
	for _, s := range ss {
		for k, v := range simpleMap {
			if strings.Contains(s, k) {
				output = append(output, v)
				continue LOOP
			}
		}
		output = append(output, s)
	}
	return output
}

// extractModelType tries to extract a meaningful model type from a
// boxname. Since models are not well normalised further cleaning work
// is likely to be needed in future. The titling type is set to English.
// If cleaning doesn't work only the first two words of the the
// description is used.
func extractModelType(s string) string {
	titleCase := cases.Title(language.English)

	// slice of regexps and replacements
	type replacement struct {
		pattern     *regexp.Regexp // case insensitive regexp
		replacement string         // replacement string, potentially with bracketed match offset
	}

	var reReplacements = []replacement{
		// remove items after "/" character
		replacement{regexp.MustCompile(`(?i)^\s*(\w.*?)/.+`), "$1"},
		// "thinkpad" is unneeded
		replacement{regexp.MustCompile(`(?i)thinkpad\s`), ""},
		// rationalise "(Gen 3)", "Gen3", "Gen 3" etc.
		replacement{regexp.MustCompile(`(?i)\(*gen\s*([0-9]+)\)*`), "Gen$1"},
	}

	titleCleaner := func(s string) string {
		for _, r := range reReplacements {
			result := r.pattern.ReplaceAllString(s, r.replacement)
			if result != s {
				s = result
			}
		}
		return s
	}

	cleaned := titleCleaner(s)
	if cleaned != s {
		return titleCase.String(cleaned)
	}
	fields := strings.Fields(s)
	if len(fields) < 3 {
		return s
	}
	return titleCase.String(strings.Join(fields[:2], " "))
}
