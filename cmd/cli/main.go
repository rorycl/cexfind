package main

import (
	"fmt"
	"os"

	"github.com/rorycl/cexfind/search"
)

var usage = `usage : %s "search term 1" ["search term 2"...]
search Cex for second hand equipment
`

func main() {

	if len(os.Args) < 2 {
		fmt.Printf(usage, os.Args[0])
		os.Exit(1)
	}

	queries := os.Args[1:]

	results, err := search.Search(queries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k := ""
	for sortedResults := range results.Iter() {
		key, box := sortedResults.Key, sortedResults.Box
		if key != k {
			fmt.Println(key)
			k = key
		}
		fmt.Printf("   %d %s %s\n", box.Price, box.Name, box.ID)
	}
}
