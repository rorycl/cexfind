// cmd provides a common query input checker used by the console and web
// commands. The checker takes the single query input string provided by
// each commend.
package cmd

import (
	"errors"
	"strings"
)

var InputTooShortErr = errors.New("each query needs to be at least 3 characters in length")
var QuerySplitChar = ";"

// QueryInputChecker checks a query input from a command and splits it
// into individual queries using QuerySplitChar, then checks the length of
// each part.
func QueryInputChecker(inQueries ...string) ([]string, error) {
	if len(inQueries) < 1 {
		return []string{}, InputTooShortErr
	}
	outputQueries := []string{}
	for _, iq := range inQueries {
		innerQueries := strings.Split(iq, QuerySplitChar)
		for _, q := range innerQueries {
			q = strings.TrimSpace(q)
			if len(q) < 3 {
				return outputQueries, InputTooShortErr
			}
			outputQueries = append(outputQueries, q)
		}
	}
	return outputQueries, nil
}
