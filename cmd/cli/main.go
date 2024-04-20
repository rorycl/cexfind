package main

import (
	"flag"
	"fmt"
	"os"

	cex "github.com/rorycl/cexfind"
)

var usage = `
a programme to search Cex/Webuy for second hand equipment

eg <programme> [-strict] -query "query 1" [-query "query 2"...]

`

// queriesType is a flag list type
type queriesType []string

// set appends a string to a queriesType
func (q *queriesType) Set(s string) error {
	*q = append(*q, s)
	return nil
}

// String is needed for flag.Var
func (q *queriesType) String() string {
	return fmt.Sprintln(*q)
}

// indirect Exit for testing
var Exit func(code int) = os.Exit

// flagGetter indirects flagGet for testing
var flagGetter func() (queriesType, bool) = flagGet

// flagGet checks the flags
func flagGet() (queriesType, bool) {

	var strict bool
	var queries queriesType

	flag.BoolVar(&strict, "strict", false, "only return items that strictly match the search terms")
	flag.Var(&queries, "query", "list of queries")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.Parse()
	if len(queries) < 1 {
		flag.Usage()
		Exit(1)
	}

	return queries, strict
}

func main() {

	queries, strict := flagGetter()

	results, err := cex.Search(queries, strict)
	if err != nil {
		fmt.Println(err)
		Exit(1)
	}

	k := ""
	for sortedResults := range results.Iter() {
		key, box := sortedResults.Key, sortedResults.Box
		if key != k {
			fmt.Printf("\n%s\n", key)
			k = key
		}
		fmt.Printf("✱ %-3d %s %s\n      %s\n", box.Price, box.Name, box.ID, cex.URLDETAIL+box.ID)
	}
}
