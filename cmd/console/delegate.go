// taken directly from the bubbletea listfancy example, with a few
// modifications

package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"golang.design/x/clipboard"
)

func init() {
	// initialise clipboard
	err := clipboard.Init()
	if err != nil {
		panic(fmt.Sprintf("clipboard could not be initialised %s", err))
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {

	d := list.NewDefaultDelegate()
	// set item height in no of lines
	d.SetHeight(3)
	// set inter-spacing item height in no of lines
	d.SetSpacing(1)

	// override colour of default SelectedTitle
	colour := "#00FE8C"
	colourDarker := "#00E17C"
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Copy().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: colourDarker, Dark: colourDarker}).
		Foreground(lipgloss.AdaptiveColor{Light: colour, Dark: colour})

	// override colour of default SelectedDesc
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Copy().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: colourDarker, Dark: colourDarker}).
		Foreground(lipgloss.AdaptiveColor{Light: colourDarker, Dark: colourDarker})

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var uri string

		if i, ok := m.SelectedItem().(bmark); ok {
			title = i.Title()
			uri = i.URI()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				// clipboard write blocks
				go setClipboard(uri)
				return m.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		/*
			remove: key.NewBinding(
				key.WithKeys("x", "backspace"),
				key.WithHelp("x", "delete"),
			),
		*/
	}
}
