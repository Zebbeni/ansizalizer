package main

import (
	"fmt"
	"github.com/Zebbeni/ansizalizer/io"
	"os"

	"github.com/Zebbeni/ansizalizer/component/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	io.InitKeyMap()
	a := app.New()
	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		fmt.Println("Run error:", err)
		os.Exit(1)
	}
}
