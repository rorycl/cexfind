package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	cex "github.com/rorycl/cexfind"
)

func TestMainFlags(t *testing.T) {

	var exit int
	Exit = func(code int) {
		exit = code
	}

	tests := []struct {
		args        []string
		exitCode    int
		isStrict    bool
		isVerbose   bool
		hasPostcode string
		numQueries  int
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
		{
			args:       []string{"prog", "-strict", "-verbose", "-query", "query 1", "-query", "query 2"},
			exitCode:   0,
			isStrict:   true,
			isVerbose:  true,
			numQueries: 2,
		},
		{
			args:        []string{"prog", "-postcode", "SW1A 0AA", "-strict", "-verbose", "-query", "query 1", "-query", "query 2"},
			exitCode:    0,
			isStrict:    true,
			isVerbose:   true,
			hasPostcode: "SW1A 0AA",
			numQueries:  2,
		},
	}

	for i, tt := range tests {

		// reset the flag environment
		exit = 0
		flag.CommandLine = flag.NewFlagSet(fmt.Sprintf("%d", i), flag.ContinueOnError)

		os.Args = tt.args

		queries, strict, postCode, verbose := flagGet()
		t.Logf("subtest %d, args %v", i, tt.args)
		t.Logf("subtest %d, strict %v postcode %v verbose %v queries %v", i, strict, postCode, verbose, queries)
		if got, want := exit, tt.exitCode; got != want {
			t.Errorf("got exit code %d expected %d", got, want)
		}
		if tt.exitCode == 1 {
			continue
		}
		if got, want := strict, tt.isStrict; got != want {
			t.Errorf("strict got %t expected %t", got, want)
		}
		if got, want := verbose, tt.isVerbose; got != want {
			t.Errorf("verbose got %t expected %t", got, want)
		}
		if got, want := postCode, tt.hasPostcode; got != want {
			t.Errorf("postCode got %s expected %s", got, want)
		}
		if got, want := len(queries), tt.numQueries; got != want {
			t.Errorf("num queries got %d expected %d", got, want)
		}
	}
}

func TestMainMain(t *testing.T) {

	tests := []struct {
		output     string
		flagGetter func() (queriesType, bool, string, bool)
	}{
		{
			output: `
Lenovo X390
✱ 175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W10/B PALSLENX39065B
      https://uk.webuy.com/product-detail?id=PALSLENX39065B
✱ 175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39061C
      https://uk.webuy.com/product-detail?id=PALSLENX39061C
✱ 190 Lenovo X390/i5-8365U/16GB Ram/240GB SSD/13"/W11/B PALSLENX390662B
      https://uk.webuy.com/product-detail?id=PALSLENX390662B
✱ 205 Lenovo X390/i5-8265U/16GB Ram/256GB SSD/13"/W11/B PALSLENX390420B
      https://uk.webuy.com/product-detail?id=PALSLENX390420B
✱ 215 Lenovo X390/i5-8365U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39078C
      https://uk.webuy.com/product-detail?id=PALSLENX39078C
✱ 360 Lenovo X390/i7-8665U/16GB Ram/512GB SSD/13"/W11/B PALSLENX39097B
      https://uk.webuy.com/product-detail?id=PALSLENX39097B
`,
			flagGetter: func() (queriesType, bool, string, bool) {
				return queriesType{"nonstrict", "nonverbose"}, false, "", false
			},
		},
		{
			output: `showing (cash price/exchange price) stores list

Lenovo X390
✱ 175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W10/B PALSLENX39065B
      (82/116) store 10
      https://uk.webuy.com/product-detail?id=PALSLENX39065B
✱ 175 Lenovo X390/i5-8265U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39061C
      (82/116) store 4
      https://uk.webuy.com/product-detail?id=PALSLENX39061C
✱ 190 Lenovo X390/i5-8365U/16GB Ram/240GB SSD/13"/W11/B PALSLENX390662B
      (89/126) store name
      https://uk.webuy.com/product-detail?id=PALSLENX390662B
✱ 205 Lenovo X390/i5-8265U/16GB Ram/256GB SSD/13"/W11/B PALSLENX390420B
      (96/136) a specific store 2
      https://uk.webuy.com/product-detail?id=PALSLENX390420B
✱ 215 Lenovo X390/i5-8365U/8GB Ram/256GB SSD/13"/W11/C PALSLENX39078C
      (101/143) store 3
      https://uk.webuy.com/product-detail?id=PALSLENX39078C
✱ 360 Lenovo X390/i7-8665U/16GB Ram/512GB SSD/13"/W11/B PALSLENX39097B
      (169/240) store 1, store 2
      https://uk.webuy.com/product-detail?id=PALSLENX39097B
`,
			flagGetter: func() (queriesType, bool, string, bool) {
				return queriesType{"nonstrict", "verbose"}, false, "", true
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {

			// set module flagGetter from test
			flagGetter = tt.flagGetter

			f, err := os.Open("../../testdata/example.json")
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
			cex.URL = ts.URL

			// https://stackoverflow.com/a/74299854
			storeStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			main()
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = storeStdout

			if diff := cmp.Diff(tt.output, string(out)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

			log.Println("\n", string(out))

		})
	}

}
