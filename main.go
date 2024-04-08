package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

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

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, m.list.KeyMap.CursorDown):
			m.Next()
			return m, cmd // return early to override default list.CursorDown()
		case key.Matches(msg, m.list.KeyMap.CursorUp):
			m.Prev()
			return m, cmd // return early to override default list.CursorUp()
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

// Next skips down to the next non empty, non heading item utilizing
// list.CursorDown under the hood for pagination logic etc
func (m *model) Next() {
	is := m.list.Items()
	for i := 1; i < 4; i++ {
		m.list.CursorDown() // utilize list.CursorDown
		idx := m.list.Index()
		thisItem := is[idx].(item)
		if thisItem.isHeading || thisItem.desc == emptyItem {
			continue
		}
		return
	}
}

// Prev skips up to the next non empty, non heading item utilizing
// list.CursorUp under the hood for pagination logic etc
func (m *model) Prev() {
	is := m.list.Items()
	for i := 1; i < 4; i++ {
		m.list.CursorUp() // utilize list.CursorUp
		idx := m.list.Index()
		thisItem := is[idx].(item)
		if thisItem.isHeading || thisItem.desc == emptyItem {
			continue
		}
		return
	}
}

const debug bool = true

func main() {
	if debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.Println("-------------------")
	}
	items := []list.Item{
		item{desc: "this is a heading", isHeading: true},
		item{desc: "this is a normal item 1"},
		item{desc: "this is a normal item 2"},
		item{desc: "this is a normal item 3 ... and some more text"},
		item{desc: emptyItem},
		item{desc: "this is another heading", isHeading: true},
		item{desc: "this is a normal item 4"},
		item{desc: "this is a normal item 5"},
		item{desc: emptyItem},
		item{desc: "this is a heading b", isHeading: true},
		item{desc: "b this is a normal item 1"},
		item{desc: "b this is a normal item 2"},
		item{desc: "b this is a normal item 3 this is a normal item 3b this is a normal ..."},
		item{desc: emptyItem},
		item{desc: "this is another heading c", isHeading: true},
		item{desc: "c this is a normal item 4"},
		item{desc: "c this is a normal item 5"},
		item{desc: emptyItem},
		item{desc: "this is a heading d", isHeading: true},
		item{desc: "d this is a normal item 1"},
		item{desc: "d this is a normal item 2"},
		item{desc: "d this is a normal item 3 this is a normal item 3."},
		item{desc: emptyItem},
		item{desc: "this is another heading e", isHeading: true},
		item{desc: "e this is a normal item 4"},
		item{desc: "e this is a normal item 5"},
	}

	m := model{list: list.New(items, NewCustomDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"
	m.list.SetFilteringEnabled(true)
	m.list.SetShowFilter(true)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
