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
	docStyle = lipgloss.NewStyle().Margin(2, 0, 0, 3)
	// top panel
	topPanelStyle = lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ff5a56")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderBottom(true).
		Padding(0, 0, 1, 0).
		Margin(1, 0, 0, 2).
		Height(5).
		Width(60).
		Background(lipgloss.Color("#000000")).
		UnsetBold()
)

type state int

const (
	listState state = iota // default
	inputState
)

const searchPrefixTpl string = "searching for \"%s\"..."

func (s state) String() string {
	return []string{"list", "input"}[s]
}

// model contains a model for the textinput model and list model,
// together with state variables
type model struct {
	// a wrapped bubbles.textinput.Model for the input
	input tiModel
	// a wrapped bubbles.list.Model for the list elements
	list liModel
	// a wrapped string (as tea.Model) for the status updates
	status status
	// flags etc.
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
			m.status = newSelection()
			m.inited = true
		}
		var t tea.Model
		// m.list.list.Select(min(1, len(m.list.list.Items())-1))
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
		m.input.Update(msg)
		return m, cmd

	// data was entered into input and needs to be percolated to status
	case inputEnterMsg:
		log.Printf("inputEnterMsg received %v", msg)
		m.status = status(fmt.Sprintf(searchPrefixTpl, msg))
		strict := false // hardcoded for now
		return m, findPerform(string(msg), strict)

	// data was entered into input and needs to be percolated to status
	case listEnterMsg:
		log.Printf("listEnterMsg received %v", msg)
		m.status = status(msg)
		return m, nil

	// perform a web search
	case findPerformMsg:
		log.Printf("findPerformMsg received %v", msg)
		items, num, err := find(msg.query, msg.strict)
		var cmd tea.Cmd
		if err != nil {
			m.status = status("Error: " + err.Error())
		} else {
			m.status = status(fmt.Sprintf("%d items found", num))
			cmd = m.list.list.SetItems(items)
			m.stateSwitch()
		}
		return m, cmd

	case tea.KeyMsg:
		log.Printf("state %s input.focus %v key %s", m.state, m.input.input.Focused(), msg.String())
		if msg.String() == "tab" {
			log.Println("at tab")
			m.stateSwitch()
			return m, func() tea.Msg { return "" } // or nil
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
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			topPanelStyle.Render(
				m.input.View(),
				m.status.View(),
			),
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
