package dithering

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ShouldClose bool
	width       int
}

func New(w int) Model {
	return Model{
		ShouldClose: false,
		width:       w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) View() string {
	return "Dithering Content"
}
