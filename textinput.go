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
)

type tiModel struct {
	input textinput.Model
}

func newTIModel() tiModel {
	t := textinput.New()
	t.Cursor.Style = tiCursorStyle
	t.CharLimit = 55
	t.Placeholder = "enter terms"
	t.PromptStyle = tiFocusedStyle
	return tiModel{
		input: t,
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

func (ti tiModel) View() string {
	var b strings.Builder
	b.WriteString(tiNormalStyle.Render("search cex"))
	b.WriteRune('\n')
	b.WriteString(ti.input.View())
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
	ti.input, cmd = ti.input.Update(msg)
	return ti, cmd
}

// enter event message
type inputEnterMsg string
