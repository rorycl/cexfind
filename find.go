// find provides the cexfind/search compoonent to the bubbletea app
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rorycl/cexfind/search"
)

// find makes cexfind/search queries and turns the results into
// list.Items as required by the bubbletea list and delegate which uses
// the following format:
//
//		items := []list.Item{
//			item{desc: "this is a heading", isHeading: true},
//			item{desc: "this is a normal item 1"},
//
// The query is received from the app as a single string with queries
// separated (potentially) by a comma. Queries are expected to each be
// at least 4 characters in length.
func find(query string, strict bool) (items []list.Item, itemNo int, err error) {

	queries := strings.Split(query, ",")
	for _, q := range queries {
		if len(q) < 4 {
			return items, 0, errors.New("queries need to be at least 4 characters in length")
		}
	}

	var results search.BoxMap
	log.Printf("  making search for %v, strict %t", queries, strict)
	results, err = search.Search(queries, strict)
	if err != nil {
		return
	}

	k := ""
	for sortedResults := range results.Iter() {
		key, box := sortedResults.Key, sortedResults.Box
		if key != k {
			if len(items) != 0 {
				// add an empty item for spacing ahead of headings,
				// except for the first heading
				items = append(items, item{desc: emptyItem})
			}
			items = append(items, item{desc: key, isHeading: true})
			k = key
		}
		items = append(items, item{
			desc: fmt.Sprintf(boxTpl, box.Price, box.Name, box.ID),
			url:  search.URLDETAIL + box.ID,
		})
		itemNo++
	}
	return
}

const boxTpl string = "£%-3d %30s %s"

// findPerformMsg is a bubbletea Cmd message for performing a find
type findPerformMsg struct {
	query  string
	strict bool
}

// findPerform wraps a findPerformMsg in a tea.Cmd for deferred
// processing. See bubbleta/tutorials/commands
func findPerform(query string, strict bool) tea.Cmd {
	return func() tea.Msg {
		return findPerformMsg{
			query:  query,
			strict: strict,
		}
	}
}
