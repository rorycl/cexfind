// This model file is the top level model file which controls the
// subsidiary inModel search input and checkbox model and the liModel
// list model. State transitions are managed through the Update function
// here and, depending on the focus state, percolated to the subsidiary
// models.

package main

import (
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// the style container for the app
	docPanelStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 0).
			Padding(0, 0, 0, 0)

	// top panel
	topPanelStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#ff982e")).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderBottom(true).
			Height(5).
			Margin(1, 2, 0, 3).
			Width(80).
			UnsetBold()

	// an arbitrary offset amount to ensure the list panel does not push
	// the other panels off the page
	topVerticalOffset = 8

	// list panel
	listPanelStyle = lipgloss.NewStyle().
			Padding(0, 0, 0, 2).
			Margin(0, 0, 0, 1)

	// help panel
	helpPanelStyle = lipgloss.NewStyle().
			Height(1).
			Padding(0, 0, 0, 0).
			Margin(0, 0, 0, 3)
)

// state indicates the main model's state
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
	state   state
	listLen int
	inited  bool

	// keys are the current key set based on the focus state, switched
	// through getKeyMap in keymap.go
	keys help.KeyMap

	// find function indirector allows for local/testing swapping of
	// functions
	finder func(query string, strict bool, postcode string) (items []list.Item, itemNo int, err error)
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

	// set find function (normally find, but can use findLocal for
	// testing
	m.finder = find
	if _, ok := os.LookupEnv("DEBUGFIND"); ok {
		m.finder = findLocal
	}

	return &m
}

// switchStylesForListing sets the list panel style padding and margins
// to deal with empty lists or the pagination marker that only appears
// at the bottom of a listing when more than one page of results are
// found.
func (m model) switchStylesForListing() {
	listPanelStyle.PaddingLeft(0)
	if m.listLen == 0 || m.list.list.Paginator.TotalPages < 1 {
		// offset padding of empty items
		listPanelStyle.PaddingLeft(2)
	}
	helpPanelStyle.MarginTop(1)
	if m.list.list.Paginator.TotalPages > 1 {
		// when listing force the help panel up
		helpPanelStyle.MarginTop(0)
	}
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
			m.status = m.status.setInputting()
		}
	case checkboxState:
		m.input.cursor = cursorBox
		m.input.Blur()
		m.keys = getKeyMap(inputKeysState)
		if withStatus {
			m.status = m.status.setCheckbox()
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
		dh, dv := docPanelStyle.GetFrameSize()
		th, tv := topPanelStyle.GetFrameSize()
		lh, lv := listPanelStyle.GetFrameSize()
		hh, hv := helpPanelStyle.GetFrameSize()

		windowWidth := msg.Width
		windowHeight := msg.Height

		w := windowWidth - dh - th - lh - hh
		h := windowHeight - dv - tv - lv - hv - topVerticalOffset

		log.Printf("window %d:%d remaining %d:%d", windowWidth, windowHeight, w, h)
		m.list.list.SetSize(w, h)
		var t tea.Model

		t, cmd = m.list.Update(msg)
		m.list = t.(liModel)
		m.input.Update(msg)
		return m, cmd

	// search data was entered into input and needs to be percolated to status
	case inputEnterMsg:
		log.Printf("inputEnterMsg received %v", msg)
		m.status = m.status.setSearching(string(msg))
		return m, findPerform(string(msg), m.input.checkbox, m.input.postcode.Value())

	// data was selected in the list view; reset the status after a
	// short wait
	case listEnterMsg:
		log.Printf("listEnterMsg received %#v", msg)
		if err := clipboard.WriteAll(msg.url); err != nil {
			m.status = m.status.setNotCopied(msg.String())
		} else {
			m.status = m.status.setCopied(msg.String())
		}
		return m, func() tea.Msg {
			time.Sleep(2500 * time.Millisecond)
			return resetListStatus{}
		}

	// perform a web search
	case findPerformMsg:
		time.Sleep(250 * time.Millisecond) // give time for status to show
		log.Printf("findPerformMsg received %v", msg)
		items, num, err := m.finder(msg.query, msg.strict, msg.postcode)
		var cmd tea.Cmd
		switch {
		case num > 0 && err != nil:
			// show the list results and change focus there, but also
			// 1. show the error for a second in the status area
			// 2. then show the normal "found" status
			m.status = status("Error: " + err.Error())
			m.listLen = num
			cmd = func() tea.Msg {
				time.Sleep(1250 * time.Millisecond)
				return resetListStatus{}
			}
			cmds = append(cmds, cmd)
			cmd = m.list.ReplaceList(items)
			cmds = append(cmds, cmd)
			cmd = m.stateSwitch(listState, false) // switch to list view
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		case err != nil:
			// the straight-forward error case
			m.status = status("Error: " + err.Error())
			emptyItem := item{}
			emptyList := []list.Item{emptyItem} // empty list
			cmd = m.list.ReplaceList(emptyList)
			cmds = append(cmds, cmd)
			cmd = m.stateSwitch(inputState, false)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		default:
			// standard, non-error list
			m.listLen = num
			m.status = m.status.setFoundItems(num)
			cmd = m.list.ReplaceList(items)
			cmds = append(cmds, cmd)
			cmd = m.stateSwitch(listState, false) // switch to list view
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

	// reset the list status
	case resetListStatus:
		m.status = m.status.setFoundItems(m.listLen)
		return m, nil

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
			cmd = m.stateSwitch(s, true)
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
	m.switchStylesForListing() // fix list panel styling if needed
	return docPanelStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			topPanelStyle.Render(
				m.input.View(),
				m.status.View(),
			),
			listPanelStyle.Render(m.list.View()),
			helpPanelStyle.Render(m.help.View(m.keys)),
		),
	)
}
