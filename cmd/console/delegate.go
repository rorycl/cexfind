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

	// A Normal description
	NormalDescription   lipgloss.Style
	SelectedDescription lipgloss.Style
	DimmedDesc          lipgloss.Style // default
	// DimmedDescription lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewCustomItemStyles returns style definitions for a default item. See
// CustomItemView for when these come into play.
func NewCustomItemStyles() (s CustomItemStyles) {

	s.Heading = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#b7b7b7", Dark: "#b7b7b7"}).
		Padding(0, 0, 0, 2).
		Bold(true)
	s.SelectedHeading = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#b7b7b7", Dark: "#b7b7b7"}).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Padding(0, 0, 0, 2).
		Bold(true)

	s.NormalDescription = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#b7b7b7", Dark: "#b7b7b7"}).
		Padding(0, 0, 0, 2)
	s.SelectedDescription = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#ff982e", Dark: "#ff982e"}).
		Padding(0, 0, 0, 1)

	return s
}

// CustomItem describes an items designed to work with CustomDelegate.
type CustomItem interface {
	list.Item
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
		isHeading bool
		desc      string
		s         = &d.Styles
	)

	if i, ok := item.(CustomItem); ok {
		desc = i.Description()
		isHeading = i.IsHeading()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width (see original // implementation)
	textwidth := uint(m.Width() - s.NormalDescription.GetPaddingLeft() - s.NormalDescription.GetPaddingRight())
	var lines []string
	for _, line := range strings.Split(desc, "\n") {
		lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
	}
	desc = strings.Join(lines, "\n")

	// Conditions
	var (
		isSelected = index == m.Index()
		isEmpty    = desc == emptyItem
	)

	if isEmpty {
		desc = ""
	}
	if isHeading {
		switch {
		case isSelected:
			desc = s.SelectedHeading.Render(desc)
		default:
			desc = s.Heading.Render(desc)
		}
	} else {
		switch {
		case isSelected:
			desc = s.SelectedDescription.Render(desc)
		default:
			desc = s.NormalDescription.Render(desc)
		}
	}

	fmt.Fprintf(w, "%s", desc)
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
