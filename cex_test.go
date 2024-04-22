package cexfind

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	results, err := Search([]string{"lenovo x390s"}, false)
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
	_, err = Search([]string{"lenovo x390st"}, true)
	if err == nil || err.Error() != "no results" {
		t.Fatalf("expected no results error, got %v", err)
	}

}

// TestBoxSort tests box sorting
func TestBoxSort(t *testing.T) {

	var toSortBoxes boxes
	toSortBoxes = append(toSortBoxes,
		[]Box{
			{"bb", "bb", "id1", 20},
			{"bc", "cc", "id2", 25},
			{"ba", "aa", "id3", 15},
			{"ab", "db", "id3", 30},
			{"ac", "dc", "id2", 35},
			{"aa", "da", "id1", 35},
			{"aa", "la", "id1", 30},
		}...,
	)

	var sortedBoxes = make(boxes, len(toSortBoxes))
	copy(sortedBoxes, toSortBoxes)
	sortedBoxes.sort()

	if diff := cmp.Diff(sortedBoxes[0], toSortBoxes[6]); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(sortedBoxes[4], toSortBoxes[2]); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

// TestBoxIDUrl checks a valid url is returned
func TestBoxIDUrl(t *testing.T) {
	b := Box{ID: "xyz"}
	if got, want := b.IDUrl(), urlDetail+b.ID; got != want {
		t.Errorf("url got %s want %s", got, want)
	}
}
