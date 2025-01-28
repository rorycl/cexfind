// A cli client to github.com/rorycl/cexfind
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rorycl/cexfind"
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
	infoStyle = color.New(color.FgMagenta).SprintFunc()
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
var flagGetter func() (queriesType, bool, string, bool) = flagGet

// flagGet checks the flags
func flagGet() (queriesType, bool, string, bool) {

	var (
		strict   bool
		queries  queriesType
		postCode string
		verbose  bool
	)

	flag.BoolVar(&strict, "strict", false, "only return items that strictly match the search terms")
	flag.Var(&queries, "query", "list of queries")
	flag.BoolVar(&verbose, "verbose", false, "show verbose output, including cash/exchange prices and stores")
	flag.StringVar(&postCode, "postcode", "", "specify postcodepostcode")

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

	return queries, strict, postCode, verbose
}

func main() {

	queries, strict, postCode, verbose := flagGetter()

	// clean queries
	queries, err := cmd.QueryInputChecker(queries...)
	if err != nil {
		fmt.Println(err)
		Exit(1)
	}

	// do search
	cex := cexfind.NewCexFind()
	results, err := cex.Search(queries, strict, postCode)
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

	if verbose || postCode != "" {
		// print header
		fmt.Print("showing (cash/exchange price) and stores list")
		if postCode != "" {
			if !cex.LocationDistancesOK() {
				fmt.Print("\nnote: distance calculations failed.")
			} else {
				fmt.Print(", distance to stores in miles.")
			}
		}
		fmt.Println("")
	}

	k := ""
	for _, box := range results {
		if box.Model != k {
			fmt.Printf("\n%s\n", box.Model)
			k = box.Model
		}
		if verbose || postCode != "" {
			info := fmt.Sprintf("(%d/%d) %s",
				box.PriceCash.IntPart(),
				box.PriceExchange.IntPart(),
				box.StoresString(80),
			)
			fmt.Printf(
				"%s %-3d %s [%s]\n      %s\n      %s\n",
				dotStyle("✱"),
				box.Price.IntPart(),
				box.Name,
				box.Category,
				// box.ID,
				urlStyle(box.IDUrl()),
				infoStyle(info),
			)
		} else {
			fmt.Printf(
				"%s %-3d %s [%s]\n      %s\n",
				dotStyle("✱"),
				box.Price.IntPart(),
				box.Name,
				box.Category,
				// box.ID,
				urlStyle(box.IDUrl()),
			)

		}
	}
}
