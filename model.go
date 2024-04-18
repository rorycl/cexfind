// This model file is the top level model file which controls the
// subsidiary inModel search input and checkbox model and the liModel
// list model. State transitions are managed through the Update function
// here and, depending on the focus state, percolated to the subsidiary
// models.

package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"
)

var (
	docStyle = lipgloss.NewStyle().Margin(2, 0, 0, 3)
	// top panel
	topPanelStyle = lipgloss.NewStyle().
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
	checkboxState
)

func (s state) String() string {
	return []string{"list", "input", "checkbox"}[s]
}

// model contains a model for the textinput model and list model,
// together with state variables
type model struct {
	// a wrapped bubbles.textinput.Model for the input
	input inModel
	// a wrapped bubbles.list.Model for the list elements
	list liModel
	// a wrapped string (as tea.Model) for the status updates
	status status
	// a standard help model
	help help.Model

	// flags etc.
	state       state
	listLen     int
	inited      bool
	clipboardOK bool

	// keys are the current key set based on the focus state, switched
	// through getKeyMap in keymap.go
	keys help.KeyMap

	// find function indirector allows for local/testing swapping of
	// functions
	finder func(query string, strict bool) (items []list.Item, itemNo int, err error)
}

// NewModel creates a new model containing the input, status and list
// models within it. Focus starts in the input model
func NewModel() *model {
	m := model{
		input: newInModel(),
		list:  newLiModel(),
	}
	m.state = inputState
	m.input.Focus()
	m.status = newSelection()
	m.inited = true
	m.listLen = 0

	// initialise the help model and related keys
	m.help = help.New()
	m.help.ShowAll = false // only show short help
	m.keys = getKeyMap(inputKeysState)

	// check if clipboard can run
	err := clipboard.Init()
	if err == nil {
		m.clipboardOK = true
	}

	// set find function (normally find, but can use findLocal for
	// testing
	// m.finder = find
	m.finder = findLocal

	return &m
}

func (m model) Init() tea.Cmd {
	return m.input.Init()
}

// stateSwitch switches state between the input, checkbox and list.
// The input and checkbox are part of the input model, while the list is
// separate. tea.Msgs (which are converted to tea.Cmds) are triggered on
// state switch to update the status area. The relevant key.KeyMap is
// selected based on the current focus area
func (m *model) stateSwitch(targetState state, withStatus bool) tea.Cmd {
	defer log.Printf("state %s input.cursor %d input.focus %v", m.state, m.input.cursor, m.input.input.Focused())
	m.state = targetState
	switch targetState {
	case inputState:
		m.input.cursor = cursorInput
		m.input.Focus()
		m.keys = getKeyMap(inputKeysState)
		if withStatus {
			m.status.setInputting()
		}
	case checkboxState:
		m.input.cursor = cursorBox
		m.input.Blur()
		m.keys = getKeyMap(inputKeysState)
		if withStatus {
			m.status.setCheckbox()
		}
	case listState:
		m.input.cursor = cursorInput
		m.input.Blur()
		m.keys = getKeyMap(listKeysState)
		m.state = listState
	}
	return nil
}

// Update is a required bubbletea function and is the programme's main
// update loop which also calls each subsidiary model's Update function
// when appropriate.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var t tea.Model
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
		m.input.Update(msg)
		return m, cmd

	// search data was entered into input and needs to be percolated to status
	case inputEnterMsg:
		log.Printf("inputEnterMsg received %v", msg)
		m.status = m.status.setSearching(string(msg))
		return m, findPerform(string(msg), m.input.checkbox)

	// data was selected in the list view
	case listEnterMsg:
		log.Printf("listEnterMsg received %#v", msg)
		if m.clipboardOK {
			clipboard.Write(clipboard.FmtText, []byte(msg.url))
			m.status = m.status.setCopied(msg.String())
		} else {
			m.status = m.status.setNotCopied(msg.String())
		}
		return m, nil

	// perform a web search
	case findPerformMsg:
		time.Sleep(250 * time.Millisecond) // give time for status to show
		log.Printf("findPerformMsg received %v", msg)
		items, num, err := m.finder(msg.query, msg.strict)
		var cmd tea.Cmd
		if err != nil {
			m.status = status("Error: " + err.Error())
			cmd = m.stateSwitch(inputState, false)
			return m, cmd
		}
		m.listLen = num
		m.status = m.status.setFoundItems(num)
		cmd = m.list.ReplaceList(items)
		cmds = append(cmds, cmd)
		cmd = m.stateSwitch(listState, false) // switch to list view
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	// Catch tab here for switching between input, checkbox and list
	// One can use
	//
	//		if msg.String() == "tab" {
	//
	// for key matching, or use key.Matches
	case tea.KeyMsg:
		// log.Printf("state %s input.focus %v key '%s'", m.state, m.input.input.Focused(), msg.String())
		if key.Matches(msg, inputKeys.Tab) {
			var s state
			var w bool = true
			switch m.state {
			case inputState:
				s = checkboxState
			case checkboxState:
				if m.listLen > 0 {
					s = listState
				} else {
					s = inputState
				}
			case listState:
				s = inputState
			}
			cmd = m.stateSwitch(s, w)
			return m, cmd
		}
	}

	// defer to input or list models depending on state
	switch m.state {
	case inputState, checkboxState:
		var t tea.Model
		t, cmd = m.input.Update(msg)
		m.input = t.(inModel)
	case listState:
		var t tea.Model
		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
	}
	return m, cmd
}

// bubbletea View; this is the main view bringing the subsidiary views
// together
func (m model) View() string {
	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			topPanelStyle.Render(
				m.input.View(),
				m.status.View(),
			),
			m.list.View(),
			m.help.View(m.keys),
		),
	)
}
