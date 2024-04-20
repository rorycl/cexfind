package search

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// testTypeExtraction attempts to extract equipment types from their
// long names
func TestTypeExtraction(t *testing.T) {

	for i, r := range []struct {
		name  string
		typer string
	}{
		{
			name:  `Lenovo T490s/i5-8365U/8GB Ram/256GB SSD/14"/W10/C PALSLENT490S72C`,
			typer: "Lenovo T490s",
		},
		{
			name:  `Lenovo T495S/Ryzen3500U/16GB Ram/256GB SSD/14"/W10/B PALSLENT495S26B`,
			typer: "Lenovo T495s",
		},
		{
			name:  `Lenovo Tab TB-X306F M10 HD Gen32GB" Iron Gray, WiFi B STABLENTBX306F32IGWB`,
			typer: "Lenovo Tab",
		},
		{
			name:  `Lenovo TAB4 TB-X304FGB" Black, WiFi C TABLESXTBX304F16GBC`,
			typer: "Lenovo Tab4",
		},
		{
			name:  `Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39061C`,
			typer: "Lenovo X390",
		},
		{
			name:  `Lenovo XT90 True Wireless In-Ear Earphones, XYZTNLTZ AA`,
			typer: "Lenovo Xt90",
		},
		{
			name:  `Lenovo Thinkpad T14S/Ryzen3500U/16GB Ram/256GB SSD/14"/W10/B PALSLENTXXXXXXB`,
			typer: "Lenovo T14s",
		},
		{
			name:  `Lenovo T14 Gen 1/i7-10610U/32GB Ram/512GB SSD/14"/MX330/W10/B PALSLENT14G178B`,
			typer: "Lenovo T14 Gen1",
		},
		{
			name:  `Lenovo T14 Gen4/i7-1355u/16GB RAM/512GB SSD/14"/W11/A PALSLENT14G4142A`,
			typer: "Lenovo T14 Gen4",
		},
		{
			name:  `Lenovo T14 (Gen3)/i5-1245U/16GB Ram/512GB SSD/14"/W11/B PALSLENT14GEN3514B`,
			typer: "Lenovo T14 Gen3",
		},
	} {
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if got, want := extractModelType(r.name), r.typer; got != want {
				t.Errorf("got %s != want %s", got, want)
			}
			t.Logf("%s <-- %s", extractModelType(r.name), r.name)
		})
	}
}

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

// TestHeadingExtract tests if extracting an h1 heading from a stream of
// bytes works
func TestHeadingExtract(t *testing.T) {

	f, err := os.Open("testdata/error.html")
	if err != nil {
		t.Fatal(err)
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	for i, tt := range []struct {
		input  []byte
		output string
	}{
		{
			input:  []byte("hi there <h1>this is some test</h1> ok"),
			output: "this is some test",
		},
		{
			input:  []byte("xyz"),
			output: "",
		},
		{
			input:  contents,
			output: "Sorry, you have been blocked",
		},
	} {
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if got, want := headingExtract(tt.input), tt.output; got != want {
				t.Errorf("got %s != want %s", got, want)
			}
			shortInput := tt.input
			if len(shortInput) > 30 {
				shortInput = slices.Concat(shortInput[:30], []byte("..."))
			}
			t.Logf("%s resulted in `%s`", string(shortInput), headingExtract(tt.input))
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

	if got, want := len(results), 1; got != want {
		t.Fatalf("expected %d type/model result, got %d", want, got)
	}
	for _, v := range results {
		if got, want := len(v), 6; got != want {
			t.Errorf("expected %d items in the result, got %d", want, got)
		}
	}

	// verbose output (use test -v)
	for k, v := range results {
		t.Log(k)
		for _, b := range v {
			t.Log("\t", b)
		}
	}

	// strict search for non-existing model
	_, err = Search([]string{"lenovo x390st"}, true)
	if err == nil || err.Error() != "no results" {
		t.Fatalf("expected no results error, got %v", err)
	}

}

// TestBoxMapIter iterates over a BoxMap container in key order and then
// by Box Model
func TestBoxMapIter(t *testing.T) {

	boxes := BoxMap{
		"b": []Box{
			{"bb", "bb", "id1", 20},
			{"bc", "cc", "id2", 25},
			{"ba", "aa", "id3", 15},
		},
		"a": []Box{
			{"ab", "db", "id3", 30},
			{"ac", "dc", "id2", 35},
			{"aa", "da", "id1", 35},
		},
	}

	all := []boxMapIter{}
	for bi := range boxes.Iter() {
		all = append(all, bi)
	}

	if diff := cmp.Diff(all[0].Box, boxes["a"][0]); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(all[3].Box, boxes["b"][0]); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

}