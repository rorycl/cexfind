// This input file contains the bubbletea "inModel" model code for the
// search text input box and associated checkbox. Although it would be
// possible to separate these into separate model files, a single model
// is used to control both.

package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// intro text etc
	inNormalStyle = lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#1ed71a"))
		Foreground(lipgloss.Color("#ff5a56"))

	// search
	inFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			PaddingTop(1)

	// search cursor style
	inCursorStyle = inFocusedStyle.Copy()

	// checkbox
	checkBoxFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#d7d7d7", Dark: "#d7d7d7"}).
				Bold(true).
				PaddingTop(1)
	checkBoxNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#ff5a56", Dark: "#ff5a56"}).
				PaddingTop(1)
)

// inCursor tracks the cursor state between the input and checkboxes
type inCursor int

const (
	cursorInput inCursor = iota
	cursorBox
)

// inModel is the main model
type inModel struct {
	input    textinput.Model
	checkbox bool
	cursor   inCursor
}

// newInModel constructs a new inModel
func newInModel() inModel {
	t := textinput.New()
	t.Cursor.Style = inCursorStyle
	t.CharLimit = 55
	t.Placeholder = "enter terms"
	t.PromptStyle = inFocusedStyle
	t.Width = 65

	// t.KeyMap = *inputKeyMap()

	return inModel{
		input:    t,
		checkbox: false,
	}
}

// Init is a bubbletea required function
func (in inModel) Init() tea.Cmd {
	return textinput.Blink
}

func (in *inModel) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	in.input, cmd = in.input.Update(msg)
	return cmd
}

func (in *inModel) Focus() {
	in.input.Focus()
}

func (in *inModel) Blur() {
	in.input.Blur()
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
		switch msg.String() {
		case "ctrl+c", "esc":
			return in, tea.Quit
		case "enter":
			return in, func() tea.Msg {
				return inputEnterMsg(in.input.Value())
			}
		}
	}

	if in.cursor == cursorBox { // focus is on checkbox
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "space", " ", "x":
				if in.checkbox {
					in.checkbox = false
				} else {
					in.checkbox = true
				}
				return in, nil
			}
		}
	}

	in.input, cmd = in.input.Update(msg)
	return in, cmd
}

// enter event message
type inputEnterMsg string

// other status update messages
type statusUpdateMsg string
