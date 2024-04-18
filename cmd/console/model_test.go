package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func TestMain(t *testing.T) {

	// does not add empty items between headings as done by "find" (see
	// find.go)
	items := []list.Item{
		item{desc: "this is a heading", isHeading: true},
		item{desc: "this is a normal item 1", url: "https://test.com/abc/a"},
		item{desc: "this is a normal item 2", url: "https://test.com/abc/b"},
		item{desc: "this is a normal item 3 ... and some more text", url: "https://test.com/abc/c"},
		item{desc: "this is another heading", isHeading: true},
		item{desc: "this is a normal item 4", url: "https://test.com/abc/d"},
		item{desc: "this is a normal item 5", url: "https://test.com/abc/e"},
		item{desc: "this is a heading b", isHeading: true},
		item{desc: "b this is a normal item 1", url: "https://test.com/abc/f"},
		item{desc: "b this is a normal item 2", url: "https://test.com/abc/g"},
		item{desc: "b this is a normal item 3 this is a normal item 3b this is a normal ...", url: "https://test.com/abc/h"},
		item{desc: "this is another heading c", isHeading: true},
		item{desc: "c this is a normal item 4", url: "https://test.com/abc/i"},
		item{desc: "c this is a normal item 5", url: "https://test.com/abc/j"},
		item{desc: "this is a heading d", isHeading: true},
		item{desc: "d this is a normal item 1", url: "https://test.com/abc/k"},
		item{desc: "d this is a normal item 2", url: "https://test.com/abc/l"},
		item{desc: "d this is a normal item 3 this is a normal item 3.", url: "https://test.com/abc/m"},
		item{desc: "this is another heading e", isHeading: true},
		item{desc: "e this is a normal item 4", url: "https://test.com/abc/n"},
		item{desc: "e this is a normal item 5", url: "https://test.com/abc/o"},
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
