package cexfind

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rorycl/cexfind/location"
	"github.com/shopspring/decimal"
)

// TestBoxInQuery tests strict query/Box.Name matches
func TestBoxInQuery(t *testing.T) {

	tests := []struct {
		box     Box
		queries []string
		result  bool
	}{
		{
			box:     Box{Name: "ABC def hij"},
			queries: []string{"xyz ntz", "hij abc"},
			result:  true,
		},
		{
			box:     Box{Name: "ABC def hij"},
			queries: []string{"xyz ntz", "hij dbc"},
			result:  false,
		},
		{
			box:     Box{Name: "abc def hij"},
			queries: []string{"HIJ ABC"},
			result:  true,
		},
		{
			box:     Box{Name: "abc def hij", Model: "lenovo"},
			queries: []string{"HIJ ABC Lenovo"},
			result:  true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if got, want := tt.box.inQuery(tt.queries), tt.result; got != want {
				t.Errorf("got %t != want %t", got, want)
			}
		})
	}
}

func TestSearch(t *testing.T) {

	f, err := os.Open("testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(contents))
	}))
	defer ts.Close()

	// overwrite global URL with test URL
	URL = ts.URL

	// non-strict search
	cex := NewCexFind()
	results, err := cex.Search([]string{"lenovo x390s"}, false, "")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(results), 6; got != want {
		t.Fatalf("expected %d box results, got %d", want, got)
	}

	// verbose output (use test -v)
	for _, v := range results {
		t.Log("\t", v)
	}

	// strict search for non-existing model
	cex = NewCexFind()
	_, err = cex.Search([]string{"lenovo x390st"}, true, "")
	if err == nil || err.Error() != "no results" {
		t.Fatalf("expected no results error, got %v", err)
	}
}

// Search for terminator search string
func TestSearchTerminator(t *testing.T) {

	f, err := os.Open("testdata/terminator.json")
	if err != nil {
		t.Fatal(err)
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(contents))
	}))
	defer ts.Close()

	// overwrite global URL with test URL
	URL = ts.URL

	// non-strict search
	cex := NewCexFind()
	results, err := cex.Search([]string{"terminator"}, false, "")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(results), 17; got != want {
		t.Fatalf("expected %d box results, got %d", want, got)
	}

	// verbose output (use test -v)
	for _, v := range results {
		t.Log("\t", v)
	}
}

// TestBoxSort tests box sorting
func TestBoxSort(t *testing.T) {

	var toSortBoxes boxes
	toSortBoxes = append(toSortBoxes,

		[]Box{
			{Model: "bb", Name: "bb", ID: "id1a", Price: decimal.NewFromInt(20), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "bc", Name: "cc", ID: "id2a", Price: decimal.NewFromInt(25), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "ba", Name: "aa", ID: "id3a", Price: decimal.NewFromInt(15), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "ab", Name: "db", ID: "id3b", Price: decimal.NewFromInt(30), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "ac", Name: "dc", ID: "id2z", Price: decimal.NewFromInt(35), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "aa", Name: "da", ID: "id1a", Price: decimal.NewFromInt(35), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)},
			{Model: "aa", Name: "la", ID: "id1b", Price: decimal.NewFromInt(30), PriceCash: decimal.NewFromInt(15), PriceExchange: decimal.NewFromInt(17)}, // 0
		}...,
	)

	var sortedBoxes = make(boxes, len(toSortBoxes))
	copy(sortedBoxes, toSortBoxes)

	sortedBoxes.sort()

	// t.Logf("\n%d: %v\n", len(toSortBoxes), toSortBoxes)
	// t.Logf("\n%d: %v\n", len(sortedBoxes), sortedBoxes)
	// t.Log(sortedBoxes)

	// compaction does not happen here
	if got, want := len(sortedBoxes), len(toSortBoxes); got != want {
		t.Errorf("expected compaction want %d items, got %d", want, got)
	}

	if diff := cmp.Diff(
		toSortBoxes[6],
		sortedBoxes[0],
		cmpopts.IgnoreFields(Box{}, "storeNames"),
	); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(
		toSortBoxes[0],
		sortedBoxes[5],
		cmpopts.IgnoreFields(Box{}, "storeNames"),
	); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestBoxStoresString(t *testing.T) {
	tests := []struct {
		length    int
		boxString []string
		want      string
	}{
		{
			length:    6,
			boxString: []string{},
			want:      "",
		},
		{
			length:    7,
			boxString: []string{"a", "b", "c"},
			want:      "a, b, c",
		},
		{
			length:    5,
			boxString: []string{"a", "b", "c"},
			want:      "a, b…",
		},
		{
			length:    14,
			boxString: []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			want:      "a, b, c, d, e…",
		},
	}
	for i, tt := range tests {
		box := Box{ID: "whatever"}
		for _, bs := range tt.boxString {
			box.Stores = append(box.Stores, location.StoreWithDistance{
				StoreName: bs,
			})
		}
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if got, want := box.StoresString(tt.length), tt.want; got != want {
				t.Errorf("got %s != want %s", got, want)
			}
		})
	}
}

// TestBoxIDUrl checks a valid url is returned
func TestBoxIDUrl(t *testing.T) {
	b := Box{ID: "xyz"}
	if got, want := b.IDUrl(), urlDetail+b.ID; got != want {
		t.Errorf("url got %s want %s", got, want)
	}
}

func TestCexInitialised(t *testing.T) {
	cex := &CexFind{
		storeDistances: location.NewStoreDistances(false),
	}
	if got, want := cex.LocationDistancesOK(), false; got != want {
		t.Errorf("got %t want %t", got, want)
	}
}
