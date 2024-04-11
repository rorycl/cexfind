/*
status shows status messages in the app, receiving and displaying
bubbletea cmd (tea.Cmd) messages to help instruct the user.
*/

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A3A3A3")).
		PaddingTop(2)
)

type status string

func newSelection() status {
	return status("add searches separated by a comma")
}

func (s status) Init() tea.Cmd {
	return nil
}

func (s status) View() string {
	return statusStyle.Render(string(s))
}

func (s status) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}
