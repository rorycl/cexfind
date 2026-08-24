package main

import (
	"fmt"
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

	m, err := NewModel("")
	if err != nil {
		t.Fatal(err)
	}
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

func TestNewModelFailures(t *testing.T) {

	tests := []struct {
		proxy string
		ok    bool
	}{
		{"socks5://127.0.0.1:8081", true},
		{"socks7://127.0.0.1:8081", false},
		{"nonsense", false},
	}

	for ii, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", ii), func(t *testing.T) {
			_, err := NewModel(tt.proxy)
			if err != nil && tt.ok {
				t.Errorf("error %s, expected none", err)
			}
			if err == nil && !tt.ok {
				t.Errorf("go no error, expected one")
			}
		})
	}
}
