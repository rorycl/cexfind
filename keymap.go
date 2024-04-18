// help and action keymaps
package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type keyState int

const (
	inputKeysState keyState = iota
	listKeysState
)

// getKeyMap returns a key.KeyMap
func getKeyMap(k keyState) help.KeyMap {
	if k == inputKeysState {
		return inputKeys
	}
	return listKeys
}

type inputKeyMap struct {
	Search   key.Binding // enter, do the search
	Tab      key.Binding // switch focus
	Selector key.Binding // select or deselect strict
	Exit     key.Binding // exit the app
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Tab, k.Selector, k.Exit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var inputKeys = inputKeyMap{
	Search: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "search"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change focus"),
	),
	Selector: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("<space>,x", "de/select strict"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("^c", "exit"),
	),
}

// KeyMap is a reduced list of default items from the default at
// "github.com/charmbracelet/bubbles/key"
type KeyMap struct {
	// common
	Enter key.Binding // select item in list
	Tab   key.Binding // switch focus
	Exit  key.Binding // exit the app
	// Keybindings used when browsing the list.
	CursorUp   key.Binding
	CursorDown key.Binding
	// keybindings such as NexTpage and PrevPage continue to be caught
	// by the default list keymap
	/*
		NextPage   key.Binding
		PrevPage   key.Binding
	*/
}

// ShortHelp
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Tab, k.Exit, k.CursorUp, k.CursorDown}
}

// FullHelp disabled
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var listKeys = KeyMap{
	// common
	Enter: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter/space", "select"),
	),
	Tab:  inputKeys.Tab,
	Exit: inputKeys.Exit,

	// browsing the list
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
}
