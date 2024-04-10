package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 0, 0, 3)
)

type model struct {
	input tiModel
	list  liModel
}

func (m model) Init() tea.Cmd {
	return m.input.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input.Update(msg)
	var t tea.Model
	t, cmd = m.list.Update(msg)
	m.list = t.(liModel)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.input.View(),
			m.list.View(),
		),
	)
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

	li := liModel{list.New(items, NewCustomDelegate(), 0, 0)}
	in := newTextInputModel()
	m := model{input: in, list: li}
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
