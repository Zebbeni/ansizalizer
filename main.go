package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
}

func main() {
	m := &model{}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Run error:", err)
		os.Exit(1)
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *model) View() string {
	return "Program"
}
