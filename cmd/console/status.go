// status shows status messages in the TUI, below the search bar. The
// status area shows messages set directly or received from bubbletea
// cmd (tea.Cmd) messages to help instruct the user.

package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A3A3A3")).
		PaddingTop(2)
)

type status string

// newSelection constructor
func newSelection() status {
	return status("add searches separated by a comma, tab to switch fields, enter to search")
}

// bubbletea Init
func (s status) Init() tea.Cmd {
	return nil
}

// bubbletea View
func (s status) View() string {
	return statusStyle.Render(string(s))
}

// bubbletea Update
func (s status) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

// status formatting when in input
func (s status) setInputting() status {
	return status("add searches separated by a comma, tab to switch fields, enter to search")
}

// status formatting when searching
func (s status) setSearching(t string) status {
	var searchPrefixTpl = "searching for \"%s\"..."
	return status(fmt.Sprintf(searchPrefixTpl, t))
}

// status formatting when in checkbox
func (s status) setCheckbox() status {
	return status("space or 'x' to select strict searching, enter to search")
}

// status on finding items
func (s status) setFoundItems(n int) status {
	i := "item"
	if n > 1 {
		i = "items"
	}
	tpl := "%d %s found. Copy an item's url by pressing enter"
	return status(fmt.Sprintf(tpl, n, i))
}

// status on copying to clipboard
func (s status) setCopied(msg string) status {
	tpl := `url for "%s%s" copied to clipboard`
	return status(fmt.Sprintf(tpl, msg, ellipsis))
}

// status on failing to copy to clipboard
func (s status) setNotCopied(msg string) status {
	tpl := `"%s%s" selected`
	return status(fmt.Sprintf(tpl, msg, ellipsis))
}
