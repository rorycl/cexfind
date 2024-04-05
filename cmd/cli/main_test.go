package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rorycl/cexfind/search"
)

func TestMainFlags(t *testing.T) {

	var exit int
	Exit = func(code int) {
		exit = code
	}

	tests := []struct {
		args       []string
		exitCode   int
		isStrict   bool
		numQueries int
	}{
		{
			args:     []string{"prog"},
			exitCode: 1,
		},
		{
			args:       []string{"prog", "-query", "query 1"},
			exitCode:   0,
			isStrict:   false,
			numQueries: 1,
		},
		{
			args:       []string{"prog", "-strict", "-query", "query 1", "-query", "query 2"},
			exitCode:   0,
			isStrict:   true,
			numQueries: 2,
		},
	}

	for i, tt := range tests {

		// reset the flag environment
		exit = 0
		flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("%d", i), flag.ContinueOnError)

		os.Args = tt.args

		queries, strict := flagGet()
		t.Logf("subtest %d, args %v", i, tt.args)
		t.Logf("subtest %d, strict %v queries %v", i, strict, queries)
		if got, want := exit, tt.exitCode; got != want {
			t.Errorf("got exit code %d expected %d", got, want)
		}
		if tt.exitCode == 1 {
			continue
		}
		if got, want := strict, tt.isStrict; got != want {
			t.Errorf("strict got %t expected %t", got, want)
		}
		if got, want := len(queries), tt.numQueries; got != want {
			t.Errorf("num queries got %d expected %d", got, want)
		}
	}
}

func TestMainMain(t *testing.T) {

	expectedOutput := `Lenovo X390
   175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39061C
   175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W10/B PALSLENX39065B
   190 Lenovo X390/i5-8365U/16GB Ram/240GB SSD/13"/W11/B PALSLENX390662B
   205 Lenovo X390/i5-8265U/16GB Ram/256GB SSD/13"/W11/B PALSLENX390420B
   215 Lenovo X390/i5-8365U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39078C
   360 Lenovo X390/i7-8665U/16GB Ram/512GB SSD/13"/W11/B PALSLENX39097B
`

	flagGetter = func() (queriesType, bool) {
		return queriesType{"query 1", "query2"}, false
	}

	f, err := os.Open("../../search/testdata/example.json")
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
	search.URL = ts.URL

	// https://stackoverflow.com/a/74299854
	storeStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = storeStdout

	if diff := cmp.Diff(expectedOutput, string(out)); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
