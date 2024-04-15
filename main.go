package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const debug bool = true

func main() {
	var logfile string
	if debug {
		logfile = "debug.log"
	} else {
		logfile = os.DevNull
	}
	f, err := tea.LogToFile(logfile, "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.Println("-------------------")

	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
