// delegate is a list delegate as set out in the bubble documentation for
// customising list items. This file is a modified copy of the bubbletea
// list-fancy example custom delegate with a simplified CustomItem
// interface and CustomDelegate that meets that interface.

// type DefaultItem interface {
// 	Item
// 	Title() string
// 	Description() string
// }

package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const (
	bullet   = "•"
	ellipsis = "…"
)

// CustomItemStyles defines styling for a default list item.
// See CustomItemView for when these come into play.
type CustomItemStyles struct {
	// A section heading; with a "First" variant with no padding
	Heading         lipgloss.Style
	SelectedHeading lipgloss.Style

	// Title and Description
	NormalTitle         lipgloss.Style
	SelectedTitle       lipgloss.Style
	NormalDescription   lipgloss.Style
	SelectedDescription lipgloss.Style
	DimmedDesc          lipgloss.Style // default
	// DimmedTitle lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewCustomItemStyles returns style definitions for a default item. See
// CustomItemView for when these come into play.
func NewCustomItemStyles() (s CustomItemStyles) {

	s.Heading = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#f4eee0", Dark: "#f4eee0"}).
		Padding(0, 0, 0, 2)
	s.SelectedHeading = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#b7b7b7", Dark: "#b7b7b7"}).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Padding(0, 0, 0, 2).
		Bold(true)

	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#dccbcb", Dark: "#dccbcb"}).
		Padding(0, 0, 0, 2)
	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Padding(0, 0, 0, 1)

	s.NormalDescription = s.NormalTitle.
		Faint(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#9c9c9c", Dark: "#9c9c9c"})
	s.SelectedDescription = s.SelectedTitle

	return s
}

// CustomItem describes an items designed to work with CustomDelegate.
type CustomItem interface {
	list.Item
	Title() string
	Description() string
	IsHeading() bool
}

// This section is largely copied from bubbles/list/defaultitem.go
//
// CustomDelegate is a standard delegate designed to work in lists. It's
// styled by CustomItemStyles.
//
// The spacing between items can be set with the SetSpacing method.
//
// Setting UpdateFunc is optional. If it's set it will be called when the
// ItemDelegate called, which is called when the list's Update function is
// invoked.
//
// Settings ShortHelpFunc and FullHelpFunc is optional. They can be set to
// include items in the list's default short and full help menus.
type CustomDelegate struct {
	Styles        CustomItemStyles
	UpdateFunc    func(tea.Msg, *list.Model) tea.Cmd
	ShortHelpFunc func() []key.Binding
	FullHelpFunc  func() [][]key.Binding
	height        int
	spacing       int
}

// NewCustomDelegate creates a new delegate with default styles.
func NewCustomDelegate() CustomDelegate {
	return CustomDelegate{
		Styles:  NewCustomItemStyles(),
		height:  1,
		spacing: 0,
	}
}

// SetHeight sets delegate's preferred height.
func (d *CustomDelegate) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d CustomDelegate) Height() int {
	return 2 // title + description
}

// SetSpacing sets the delegate's spacing.
func (d *CustomDelegate) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d CustomDelegate) Spacing() int {
	return d.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d CustomDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

// Render prints an item. Note filtering not used
func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		isHeading   bool
		title       string
		description string
		s           = &d.Styles
	)

	if i, ok := item.(CustomItem); ok {
		title = i.Title()
		description = i.Description()
		isHeading = i.IsHeading()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width (see original implementation)
	// nolint
	textwidth := uint(m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	var lines []string
	for _, line := range strings.Split(title, "\n") {
		lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
	}
	title = strings.Join(lines, "\n")

	// Conditions
	var (
		isSelected = index == m.Index()
		isEmpty    = title == emptyItem
	)

	// fixme (set heading)
	if isEmpty {
		title = ""
	}
	if isHeading {
		switch {
		// swap title to description to make space above
		case isSelected:
			description = s.SelectedHeading.Render(title)
		default:
			description = s.Heading.Render(title)
		}
		title = ""
	} else {
		switch {
		case isSelected:
			title = s.SelectedTitle.Render(title)
			description = s.SelectedDescription.Render(description)
		default:
			title = s.NormalTitle.Render(title)
			description = s.NormalDescription.Render(description)
		}
	}

	fmt.Fprintf(w, "%s", title)
	// if !isHeading && !isEmpty {
	fmt.Fprintf(w, "\n%s", description)
	// }
}

// ShortHelp returns the delegate's short help.
func (d CustomDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d CustomDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
