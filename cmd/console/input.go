// This input file contains the bubbletea "inModel" model code for the
// search text input box and associated checkbox. Although it would be
// possible to separate these into separate model files, a single model
// is used to control both.

package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// intro text etc
	inNormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff982e"))

	// search
	inFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			PaddingTop(1)

	// search cursor style
	inCursorStyle = inFocusedStyle

	// postcode
	postcodeFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff982e")).
				PaddingTop(1)

	postcodeNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff982e")).
				PaddingTop(1)

	// checkbox
	checkBoxFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#d7d7d7", Dark: "#d7d7d7"}).
				Bold(true).
				MarginTop(1)
	checkBoxNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
				MarginTop(1)
)

// inCursor tracks the cursor state between the input and checkboxes
type inCursor int

const (
	cursorInput inCursor = iota
	cursorPostcode
	cursorBox
)

// inModel is the main model
type inModel struct {
	input    textinput.Model
	checkbox bool
	postcode textinput.Model
	cursor   inCursor
}

// newInModel constructs a new inModel input model
func newInModel() inModel {
	t := textinput.New()
	t.Cursor.Style = inCursorStyle
	t.CharLimit = 70
	t.Placeholder = "enter terms"
	t.PromptStyle = inFocusedStyle
	t.Width = 60

	p := textinput.New()
	p.Cursor.Style = postcodeNormalStyle
	p.CharLimit = 8
	p.Placeholder = "postcode"
	p.PromptStyle = postcodeFocusedStyle
	p.Width = 12

	return inModel{
		input:    t,
		postcode: p,
		checkbox: false,
	}
}

// Init is a bubbletea required function
func (in inModel) Init() tea.Cmd {
	return textinput.Blink
}

// checkBoxAsString renders the checkbox as a string depending on its
// state and the selection status
func (in *inModel) checkBoxAsString() string {
	switch {
	case in.cursor == cursorBox && in.checkbox:
		return checkBoxFocusedStyle.Render("strict [x]")
	case in.cursor == cursorBox && !in.checkbox:
		return checkBoxFocusedStyle.Render("strict [ ]")
	case in.cursor != cursorBox && in.checkbox:
		return checkBoxNormalStyle.Render("strict [x]")
	}
	return checkBoxNormalStyle.Render("strict [ ]")
}

// View is the bubbletea View function which renders the top panel of
// the TUI, containing both the search bar and "strict" checkbox.
func (in inModel) View() string {
	var b strings.Builder
	b.WriteString(inNormalStyle.Render("search cex"))
	b.WriteRune('\n')
	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			in.input.View(),
			in.postcode.View(),
			in.checkBoxAsString(),
		),
	)
	return b.String()
}

// Update is a bubbletea required function
func (in inModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, inputKeys.Exit):
			return in, tea.Quit
		case key.Matches(msg, inputKeys.Search):
			return in, func() tea.Msg {
				return inputEnterMsg(in.input.Value())
			}
		}
	}
	// selector in cursor focus area
	if in.cursor == cursorBox {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, inputKeys.Selector):
				if in.checkbox {
					in.checkbox = false
				} else {
					in.checkbox = true
				}
				return in, nil
			}
		}
	}
	cmds := []tea.Cmd{}
	in.input, cmd = in.input.Update(msg)
	cmds = append(cmds, cmd)
	in.postcode, cmd = in.postcode.Update(msg)
	cmds = append(cmds, cmd)
	return in, tea.Batch(cmds...)
}

// enter event message
type inputEnterMsg string

// reset the list status (typically after a short-lived status message)
type resetListStatus struct{}
