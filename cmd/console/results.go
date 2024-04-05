package main

import (
	"github.com/rorycl/cexfind/search"
)

type modelResults struct {
	model string
	box   search.Box
}

// FilterValue allows a Bookmark to meet the bubbletea Item interface
func (mr modelResults) FilterValue() string {
	return mr.box.Name
}

// Title returns the title for bubbletea
func (mr modelResults) Title() string {
	return mr.box.Model
}

// URI returns the title for bubbletea
func (mr modelResults) URI() string {
	return search.URLDETAIL + mr.box.URI
}

// Description returns the description for bubbletea
func (mr modelResults) Description() string {
	return mr.box.Name
}

// Bmarks is a slice of Bookmark
type Bmarks []bmark

func getBmarks(path string) (Bmarks, error) {
	var bmarks Bmarks
	bookmarks, err := bookmark.ExtractBookmarks(path)
	if err != nil {
		return bmarks, err
	}
	for _, b := range bookmarks {
		bmarks = append(bmarks, bmark{b.Title, b.URI, b.Tags})
	}
	return bmarks, nil
}
