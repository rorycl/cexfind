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

type state int

const (
	listState state = iota // default
	inputState
)

func (s state) String() string {
	return []string{"list", "input"}[s]
}

// model contains a model for the textinput model and list model,
// together with state variables
type model struct {
	input  tiModel
	list   liModel
	state  state
	inited bool
}

func (m model) Init() tea.Cmd {
	return m.input.Init()
}

// stateSwitch switches state between the input and list panels
func (m *model) stateSwitch() {
	switch m.state {
	case inputState:
		m.state = listState
		m.input.Blur()
	default:
		m.state = inputState
		m.input.Focus()
	}
	log.Printf("state %d %s input.focus %v", m.state, m.state, m.input.input.Focused())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Println("at window resize")
		// initialise
		if !m.inited {
			m.stateSwitch()
			m.inited = true
		}
		var t tea.Model
		// m.list.list.Select(min(1, len(m.list.list.Items())-1))
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
		m.input.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		log.Printf("state %s input.focus %v key %s", m.state, m.input.input.Focused(), msg.String())
		if msg.String() == "]" {
			log.Println("at ]")
			m.stateSwitch()
		}
	}

	// defer to input or list models
	switch m.state {
	case inputState:
		var t tea.Model
		t, cmd = m.input.Update(msg)
		m.input = t.(tiModel)
	case listState:
		var t tea.Model
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
	}
	log.Printf("state %s selection %s", m.state, m.input.selection)
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
	in := newTIModel()
	m := model{input: in, list: li}
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
