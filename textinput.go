package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// intro text etc
	tiNormalStyle = lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#1ed71a"))
		Foreground(lipgloss.Color("#e7e223"))
	// search
	tiFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			PaddingTop(1)

	// search cursor style
	tiCursorStyle = tiFocusedStyle.Copy()
	// checkbox
	checkBoxFocusedStyle = tiFocusedStyle.Copy()
	checkBoxNormalStyle  = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#e7e223", Dark: "#e7e223"}).
				PaddingTop(1)
)

type tiCursor int

const (
	cursorInput tiCursor = iota
	cursorBox
)

type tiModel struct {
	input    textinput.Model
	checkbox bool
	cursor   tiCursor
}

func newTIModel() tiModel {
	t := textinput.New()
	t.Cursor.Style = tiCursorStyle
	t.CharLimit = 55
	t.Placeholder = "enter terms"
	t.PromptStyle = tiFocusedStyle
	t.Width = 65

	return tiModel{
		input:    t,
		checkbox: false,
	}
}

func (ti tiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (ti *tiModel) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	ti.input, cmd = ti.input.Update(msg)
	return cmd
}

func (ti *tiModel) Focus() {
	ti.input.Focus()
}

func (ti *tiModel) Blur() {
	ti.input.Blur()
}

func (ti *tiModel) checkBoxAsString() string {
	// log.Println("-> checkBoxAsString with checkbox set to ", ti.checkbox)
	switch {
	case ti.cursor == cursorBox && ti.checkbox:
		return checkBoxFocusedStyle.Render("strict [x]")
	case ti.cursor == cursorBox && !ti.checkbox:
		return checkBoxFocusedStyle.Render("strict [ ]")
	case ti.cursor != cursorBox && ti.checkbox:
		return checkBoxNormalStyle.Render("strict [x]")
	}
	return checkBoxNormalStyle.Render("strict [ ]")
}

func (ti tiModel) View() string {
	var b strings.Builder
	b.WriteString(tiNormalStyle.Render("search cex"))
	b.WriteRune('\n')
	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			ti.input.View(),
			ti.checkBoxAsString(),
		),
	)
	return b.String()
}

func (ti tiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return ti, tea.Quit
		case "enter":
			return ti, func() tea.Msg {
				return inputEnterMsg(ti.input.Value())
			}
		}
	}

	if ti.cursor == cursorBox { // focus is on checkbox
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "space", " ", "x":
				if ti.checkbox {
					ti.checkbox = false
				} else {
					ti.checkbox = true
				}
				return ti, nil
			}
		}
	}

	ti.input, cmd = ti.input.Update(msg)
	return ti, cmd
}

// enter event message
type inputEnterMsg string

// other status update messages
type statusUpdateMsg string
