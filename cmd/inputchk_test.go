package cmd

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindLocal(t *testing.T) {

	tests := []struct {
		input  string
		output []string
		err    error
	}{
		{
			input:  "",
			output: []string{},
			err:    InputTooShortErr,
		},
		{
			input:  "one",
			output: []string{"one"},
			err:    nil,
		},
		{
			input:  "one,two",
			output: []string{"one,two"},
			err:    nil,
		},
		{
			input:  "one;two",
			output: []string{"one", "two"},
			err:    nil,
		},
		{
			// second input too short
			input:  "one;tw",
			output: []string{"one"}, // fail
			err:    InputTooShortErr,
		},
		{
			input:  `14" laptop`,
			output: []string{`14" laptop`},
			err:    nil,
		},
		{
			input:  `   14" laptop   `,
			output: []string{`14" laptop`},
			err:    nil,
		},
		{
			input:  `   14" laptop   ; xy z `,
			output: []string{`14" laptop`, "xy z"},
			err:    nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			queries, err := QueryInputChecker(tt.input)
			if err != tt.err {
				t.Fatalf("got err %v expected err %v", err, tt.err)
			}
			if diff := cmp.Diff(tt.output, queries); diff != "" {
				t.Errorf("results unexpected %s", diff)
			}
		})
	}
}
