package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/app"
	"github.com/Zebbeni/ansizalizer/debug"
	"github.com/Zebbeni/ansizalizer/event"
)

func init() {
	event.InitKeyMap()
}

func main() {
	debugFlag := flag.Bool("debug", false, "enable debug logging to console.log and debug.log")
	flag.Parse()

	if *debugFlag {
		debug.Enable()
	}

	m := app.New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
