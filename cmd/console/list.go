// The list file contains the code for the list components of the code,
// managed through the liModel which contains a bubbles/list component.
// List items are dealt with through the list delegate CustomDelegate.

package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// an item is a list item, meeting the list.Item interface (which
// requires a
//
//	FilterValue() string
//
// function
type item struct {
	title       string // a rendered title
	description string // a rendered description
	isHeading   bool
	url         string // the url to see this item
}

// emptyItem is s special "empty" item to provide padding between item
// headings
const emptyItem = "-empty-"

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) IsHeading() bool     { return i.isHeading }
func (i item) FilterValue() string { return i.title }

type liModel struct {
	list list.Model
}

// newLiModel create a new liModel with the relevant delegate. The 0, 0
// arguments to list.New are for width and height
func newLiModel() liModel {
	li := liModel{
		list: list.New([]list.Item{}, NewCustomDelegate(), 0, 0),
	}
	li.list.SetShowTitle(false)
	li.list.SetShowStatusBar(false)
	li.list.InfiniteScrolling = false
	li.list.SetShowHelp(false) // help is customised in main model
	li.list.SetShowPagination(true)

	return li
}

// bubbletea Init
func (li liModel) Init() tea.Cmd {
	return nil
}

// bubbletea Update
func (li liModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, listKeys.Exit):
			return li, tea.Quit
		case key.Matches(msg, listKeys.CursorDown):
			li.Next()
			return li, cmd // return early to override default list.CursorDown()
		case key.Matches(msg, listKeys.CursorUp):
			li.Prev()
			return li, cmd // return early to override default list.CursorUp()
		case key.Matches(msg, listKeys.Enter):
			i := li.list.SelectedItem().(item)
			cmd = func() tea.Msg {
				return listEnterMsg{
					title: i.title,
					url:   i.url,
				}
			}
			cmds = append(cmds, cmd)
		}

		/* window size calculations are done in the main model
		case tea.WindowSizeMsg:
			...
		*/
	}

	li.list, cmd = li.list.Update(msg)
	cmds = append(cmds, cmd)
	return li, tea.Batch(cmds...)
}

// View is a bubbletea required function and renders the list component
// of the TUI window
func (li liModel) View() string {
	return li.list.View()
}

// Next skips down to the next non empty, non heading item utilizing
// list.CursorDown under the hood for pagination logic etc
func (li *liModel) Next() {
	for i := 1; i < 4; i++ {
		li.list.CursorDown() // utilize list.CursorDown
		thisItem := li.list.SelectedItem().(item)
		if thisItem.isHeading || thisItem.title == emptyItem {
			continue
		}
		return
	}
}

// Prev skips up to the next non empty, non heading item utilizing
// list.CursorUp under the hood for pagination logic etc
func (li *liModel) Prev() {
	for i := 1; i < 4; i++ {
		li.list.CursorUp() // utilize list.CursorUp
		thisItem := li.list.SelectedItem().(item)
		if thisItem.isHeading || thisItem.title == emptyItem {
			continue
		}
		return
	}
}

// ReplaceList replaces the items in the list and sets the Index
// appropriately
func (li *liModel) ReplaceList(items []list.Item) tea.Cmd {
	var cmd = li.list.SetItems(items)
	if li.list.Index() != 0 {
		li.list.Select(0)
	}
	// continue to the first non-heading, non-empty item
	thisItem := li.list.SelectedItem().(item)
	if thisItem.isHeading || thisItem.title == emptyItem {
		li.Next()
	}
	return cmd
}

// enter event message
type listEnterMsg struct {
	title string
	url   string
}

// string representation of a listEnterMsg is used for status
func (l listEnterMsg) String() string {
	trimmedDesc := l.title
	fields := strings.Fields(l.title)
	if len(fields) > 2 {
		trimmedDesc = strings.Join(fields[1:len(fields)-1], " ")
	}
	if len(trimmedDesc) > 40 {
		trimmedDesc = trimmedDesc[:40]
	}
	return trimmedDesc
}
