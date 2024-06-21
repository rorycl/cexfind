// A cli client to github.com/rorycl/cexfind
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	cex "github.com/rorycl/cexfind"
	"github.com/rorycl/cexfind/cmd"
)

var usage = `
a cli programme to search Cex/Webuy for second hand equipment

eg <programme> [-strict] -query "query 1" [-query "query 2"...]

`

// styles
var (
	urlStyle  = color.New(color.FgCyan).SprintFunc()
	dotStyle  = color.New(color.FgCyan).SprintFunc()
	infoStyle = color.New(color.FgHiWhite).SprintFunc()
)

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
var flagGetter func() (queriesType, bool, bool) = flagGet

// flagGet checks the flags
func flagGet() (queriesType, bool, bool) {

	var strict bool
	var queries queriesType
	var verbose bool

	flag.BoolVar(&strict, "strict", false, "only return items that strictly match the search terms")
	flag.Var(&queries, "query", "list of queries")
	flag.BoolVar(&verbose, "verbose", false, "show verbose output, including cash/exchange prices and stores")

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

	return queries, strict, verbose
}

func main() {

	queries, strict, verbose := flagGetter()

	// clean queries
	queries, err := cmd.QueryInputChecker(queries...)
	if err != nil {
		fmt.Println(err)
		Exit(1)
	}

	// do search
	results, err := cex.Search(queries, strict)
	switch {
	case err != nil && len(results) > 0:
		fmt.Println(err)
		// continue to show the list
	case err != nil:
		fmt.Println(err)
		Exit(1)
	default:
		// show the list
	}

	if verbose {
		// print header
		fmt.Println("showing (cash price/exchange price) stores list")
	}

	k := ""
	for _, box := range results {
		if box.Model != k {
			fmt.Printf("\n%s\n", box.Model)
			k = box.Model
		}
		if verbose {
			info := fmt.Sprintf("      (%d/%d) %s",
				box.PriceCash,
				box.PriceExchange,
				box.StoresString(),
			)
			fmt.Printf(
				"%s %-3d %s %s\n%s\n      %s\n",
				dotStyle("✱"),
				box.Price,
				box.Name,
				box.ID,
				infoStyle(info),
				urlStyle(box.IDUrl()),
			)
		} else {
			fmt.Printf(
				"%s %-3d %s %s\n      %s\n",
				dotStyle("✱"),
				box.Price,
				box.Name,
				box.ID,
				urlStyle(box.IDUrl()),
			)

		}
	}
}
