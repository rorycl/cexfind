package search

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
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
	} {
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if got, want := extractModelType(r.name), r.typer; got != want {
				t.Errorf("got %s != want %s", got, want)
			}
			t.Logf("%s <-- %s", extractModelType(r.name), r.name)
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

	results, err := Search([]string{"lenovo x390s"})
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
}
