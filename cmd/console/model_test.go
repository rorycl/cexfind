package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func TestMain(t *testing.T) {

	// does not add empty items between headings as done by "find" (see
	// find.go)
	items := []list.Item{
		item{title: "this is a heading", isHeading: true},
		item{title: "this is a normal item 1", url: "https://test.com/abc/a"},
		item{title: "this is a normal item 2", url: "https://test.com/abc/b"},
		item{title: "this is a normal item 3 ... and some more text", url: "https://test.com/abc/c"},
		item{title: "this is another heading", isHeading: true},
		item{title: "this is a normal item 4", url: "https://test.com/abc/d"},
		item{title: "this is a normal item 5", url: "https://test.com/abc/e"},
		item{title: "this is a heading b", isHeading: true},
		item{title: "b this is a normal item 1", url: "https://test.com/abc/f"},
		item{title: "b this is a normal item 2", url: "https://test.com/abc/g"},
		item{title: "b this is a normal item 3 this is a normal item 3b this is a normal ...", url: "https://test.com/abc/h"},
		item{title: "this is another heading c", isHeading: true},
		item{title: "c this is a normal item 4", url: "https://test.com/abc/i"},
		item{title: "c this is a normal item 5", url: "https://test.com/abc/j"},
		item{title: "this is a heading d", isHeading: true},
		item{title: "d this is a normal item 1", url: "https://test.com/abc/k"},
		item{title: "d this is a normal item 2", url: "https://test.com/abc/l"},
		item{title: "d this is a normal item 3 this is a normal item 3.", url: "https://test.com/abc/m"},
		item{title: "this is another heading e", isHeading: true},
		item{title: "e this is a normal item 4", url: "https://test.com/abc/n"},
		item{title: "e this is a normal item 5", url: "https://test.com/abc/o"},
	}

	m := NewModel()
	m.list.ReplaceList(items)

	if got, want := len(m.list.list.Items()), 21; got != want {
		t.Errorf("list length got %d want %d", got, want)
	}
	if got, want := m.input.cursor, cursorInput; got != want {
		t.Errorf("input cursor got %d want %d", got, want)
	}
	if got, want := m.input.checkbox, false; got != want {
		t.Errorf("checkbox set to %t want %t", got, want)
	}

}
