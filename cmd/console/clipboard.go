package main

import (
	"golang.design/x/clipboard"
)

// run setClipboard in a go routine to allow closing of out-of-date
// goroutines.
func setClipboard(text string) {
	done := clipboard.Write(clipboard.FmtText, []byte(text))
	select {
	case <-done:
		return
	}
}
