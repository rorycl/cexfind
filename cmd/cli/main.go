package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/rorycl/cexfind/search"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("please provide at least one query to search for!")
		os.Exit(1)
	}
	queries := os.Args[1:]

	allResults, err := search.Search(queries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sort for output
	keys := []string{}
	for k := range allResults {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		v := allResults[k]
		fmt.Printf("%s (%d)\n", k, len(v))
		v.Sort("Price")
		for _, iv := range v {
			fmt.Printf("   %d %s %s\n", iv.Price, iv.Name, iv.ID)
		}
	}
}
