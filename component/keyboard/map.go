package keyboard

import (
	"github.com/charmbracelet/bubbles/key"
)

type Map struct {
	Quit  key.Binding
	Nav   key.Binding
	Enter key.Binding
}

func InitMap() *Map {
	return &Map{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Nav: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑↓", "navigate"),
		),
		Enter: key.NewBinding(
			key.WithKeys("return", "enter"),
			key.WithHelp("↲", "select"),
		),
	}
}

func (k Map) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Nav, k.Enter}
}

func (k Map) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit, k.Nav, k.Enter}}
}
