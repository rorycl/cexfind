// A bubbletea TUI console client to github.com/rorycl/cexfind
// If provided with a LOGFILE env var, the program will log to it.
// If provided with a PROXY env var, the http client will route through this proxy
// (typically a local socks5 address, such as socks5://127.0.0.1:8081).
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// setup logging
	var logfile string
	if _, ok := os.LookupEnv("LOGFILE"); ok {
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

	// initialise model and cex component
	m, err := NewModel(os.Getenv("PROXY"))
	if err != nil {
		fmt.Printf("cex model initialisation error: %s\n", err)
		os.Exit(1)
	}

	// initialise bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
