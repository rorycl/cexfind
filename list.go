package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listPanel = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	// arbitrary spacing offsets for the list panel sizing
	arbitraryVerticalOffset, arbitraryHorizantalOffset = 5, 7
)

// emptyItem is s special "empty" item to provide padding between item
// headings
const emptyItem = "-empty-"

type item struct {
	desc      string
	isHeading bool
}

func (i item) Description() string { return i.desc }
func (i item) IsHeading() bool     { return i.isHeading }
func (i item) FilterValue() string { return i.desc }

type liModel struct {
	list list.Model
}

func (li liModel) Init() tea.Cmd {
	return nil
}

func (li liModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// move to Init?
	li.list.SetShowTitle(false)
	li.list.SetShowStatusBar(false)
	li.list.InfiniteScrolling = true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return li, tea.Quit
		case key.Matches(msg, li.list.KeyMap.CursorDown):
			li.Next()
			return li, cmd // return early to override default list.CursorDown()
		case key.Matches(msg, li.list.KeyMap.CursorUp):
			li.Prev()
			return li, cmd // return early to override default list.CursorUp()
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		li.list.SetSize(
			msg.Width-h-arbitraryVerticalOffset,
			msg.Height-v-arbitraryHorizantalOffset,
		)
	}

	li.list, cmd = li.list.Update(msg)
	return li, cmd
}

func (li liModel) View() string {
	return listPanel.Render(li.list.View())
}

// Next skips down to the next non empty, non heading item utilizing
// list.CursorDown under the hood for pagination logic etc
func (li *liModel) Next() {
	is := li.list.Items()
	for i := 1; i < 4; i++ {
		li.list.CursorDown() // utilize list.CursorDown
		idx := li.list.Index()
		thisItem := is[idx].(item)
		if thisItem.isHeading || thisItem.desc == emptyItem {
			continue
		}
		return
	}
}

// Prev skips up to the next non empty, non heading item utilizing
// list.CursorUp under the hood for pagination logic etc
func (li *liModel) Prev() {
	is := li.list.Items()
	for i := 1; i < 4; i++ {
		li.list.CursorUp() // utilize list.CursorUp
		idx := li.list.Index()
		thisItem := is[idx].(item)
		if thisItem.isHeading || thisItem.desc == emptyItem {
			continue
		}
		return
	}
}
