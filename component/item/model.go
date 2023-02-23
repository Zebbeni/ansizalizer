package item

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Should style live on the button itself or be applied by the rendering parent?
// also, should we have separate states for selected vs active vs hovered? Or a couple of those?

type OnSelect func()

type Model struct {
	name     string
	OnSelect OnSelect
}

func New(name string, onSelect OnSelect) Model {
	return Model{name: name, OnSelect: onSelect}
}

func (b *Model) Init() tea.Cmd {
	return nil
}

func (b *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *Model) View() string {
	return b.name
}
