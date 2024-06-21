// find provides the cexfind/search compoonent to the bubbletea app

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	cex "github.com/rorycl/cexfind"
	"github.com/rorycl/cexfind/cmd"
)

// emptyItemOn defines if "empty items" should be added to make space
// between items in the list
const emptyItemOn bool = false

// find makes cexfind.Search queries and turns the results into
// list.Items as required by the bubbletea list and delegate which uses
// the following format:
//
//	items := []list.Item{
//		item{title: "this is a heading", isHeading: true},
//		item{title: "this is a normal item 1", description: "..."},
//
// The query is received from the app as a single string with queries
// separated (potentially) by a semicolon. Queries are expected to each
// be at least 3 characters in length.
//
// The itemNo number of items is included since heading and empty items
// are introduced for formatting reasons. itemNo counts the number of
// valid items returned.
//
// Note that returning an error does not necessarily indicate a complete
// search failure since more than one query may have been made. In the
// case of an error for one query and results from others, for example,
// both an error and items will be returned.
func find(query string, strict bool) (items []list.Item, itemNo int, err error) {

	queries, err := cmd.QueryInputChecker(query)
	if err != nil {
		return items, 0, err
	}

	var results []cex.Box
	log.Printf("  making search for %v, strict %t", queries, strict)

	// note that err does not cause a failure
	results, err = cex.Search(queries, strict)

	log.Printf("results %#v\nerr %v", results, err)
	if results == nil {
		err = errors.New("no results found")
		return
	}
	itemNo = len(results)

	latestModel := ""
	for _, box := range results {
		if box.Model != latestModel {
			if len(items) != 0 && emptyItemOn {
				// add an empty item for spacing ahead of headings,
				// except for the first heading
				items = append(items, item{title: emptyItem})
			}
			// make heading
			items = append(items, item{title: box.Model, isHeading: true})
			latestModel = box.Model
		}
		// add standard item
		items = append(items, item{
			title:       fmt.Sprintf(boxTitleTpl, box.Price, box.Name),
			description: fmt.Sprintf(boxDescriptionTpl, box.PriceCash, box.PriceExchange, box.StoresString()),
			url:         box.IDUrl(),
		})
	}
	return
}

// findLocal simple returns the example list in list_example.go
func findLocal(query string, strict bool) (items []list.Item, itemNo int, err error) {
	_, err = cmd.QueryInputChecker(query)
	if err != nil {
		return items, 0, err
	}
	return theseExampleItems, 15, nil
}

const boxTitleTpl string = "£%-3d %s"
const boxDescriptionTpl string = "     (£%d/£%d) %-40s"

// findPerformMsg is a bubbletea Cmd message for performing a find
type findPerformMsg struct {
	query  string
	strict bool
}

// findPerform wraps a findPerformMsg in a tea.Cmd for deferred
// processing. See bubbleta/tutorials/commands
func findPerform(query string, strict bool) tea.Cmd {
	log.Println("   query triggered", query, strict)
	return func() tea.Msg {
		return findPerformMsg{
			query:  query,
			strict: strict,
		}
	}
}
