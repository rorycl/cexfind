package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"
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
		Width(80).
		Background(lipgloss.Color("#000000")).
		UnsetBold()
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
	// a wrapped bubbles.textinput.Model for the input
	input tiModel
	// a wrapped bubbles.list.Model for the list elements
	list liModel
	// a wrapped string (as tea.Model) for the status updates
	status status
	// flags etc.
	state       state
	inited      bool
	clipboardOK bool
}

// NewModel creates a new model containing the input, status and list
// models within it. Focus starts in the input model
func NewModel() *model {
	m := model{
		input: newTIModel(),
		list:  newLiModel(),
	}
	m.state = inputState
	m.input.Focus()
	m.status = newSelection()
	m.inited = true

	// check if clipboard can run
	err := clipboard.Init()
	if err == nil {
		m.clipboardOK = true
	}

	return &m
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
		var t tea.Model
		// m.list.list.Select(min(1, len(m.list.list.Items())-1))
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
		m.input.Update(msg)
		return m, cmd

	// data was entered into input and needs to be percolated to status
	case inputEnterMsg:
		log.Printf("inputEnterMsg received %v", msg)
		m.status = m.status.setSearching(string(msg))
		strict := false // hardcoded for now
		return m, findPerform(string(msg), strict)

	// data was selected in the list view
	case listEnterMsg:
		log.Printf("listEnterMsg received %#v", msg)
		if m.clipboardOK {
			clipboard.Write(clipboard.FmtText, []byte(msg.url))
			m.status = "url for \"" + status(msg.String()) + "\" copied to clipboard"
		} else {
			m.status = "you selected \"" + status(msg.String()) + "\""
		}
		return m, nil

	// perform a web search
	case findPerformMsg:
		time.Sleep(200 * time.Millisecond) // give time for messages to arrive
		log.Printf("findPerformMsg received %v", msg)
		items, num, err := find(msg.query, msg.strict)
		var cmd tea.Cmd
		if err != nil {
			m.status = status("Error: " + err.Error())
		} else {
			m.status = status(fmt.Sprintf("%d items found", num))
			cmd = m.list.ReplaceList(items)
			m.stateSwitch() // switch to list view
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

// bubbletea View
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
