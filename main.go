package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const debug bool = true

func main() {
	if debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.Println("-------------------")
	}
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

	/*
		li := liModel{list.New(items, NewCustomDelegate(), 0, 0)}
		in := newTIModel()
		m := model{input: in, list: li}
	*/
	m := NewModel()
	if debug {
		log.Println(items)
	}
	// m.list.ReplaceList(items)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
