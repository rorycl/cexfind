package cexfind

import (
	"fmt"
	"io"
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
