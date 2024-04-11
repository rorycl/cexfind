package main

import (
	"log"
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
	// the selection area
	tiSelectionStyle = lipgloss.NewStyle().
		//		Background(lipgloss.Color("#dbdbdb")).
		Foreground(lipgloss.Color("#A3A3A3")).
		PaddingTop(1)
	// panel
	tiPanelStyle = lipgloss.NewStyle().
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

type tiModel struct {
	input     textinput.Model
	selection string
}

func newTIModel() tiModel {
	t := textinput.New()
	t.Cursor.Style = tiCursorStyle
	t.CharLimit = 55
	t.Placeholder = "enter terms"
	t.PromptStyle = tiFocusedStyle
	return tiModel{
		input:     t,
		selection: "add searches separated by a comma",
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

func (ti *tiModel) updateSelection(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if ti.input.Value() != "" {
		ti.selection = ti.input.Value()
	}
	return cmd // empty
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
	b.WriteRune('\n')
	b.WriteString(tiSelectionStyle.Render(ti.selection))
	return tiPanelStyle.Render(b.String())
}

func (ti tiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("got into textinput key press")
		switch msg.String() {
		case "ctrl+c", "esc":
			return ti, tea.Quit
		case "enter":
			ti.updateSelection(msg)
			// return ti, nil
		}
	}
	log.Println("got into texinput key press fallthrough")
	ti.input, cmd = ti.input.Update(msg)
	return ti, cmd
}
