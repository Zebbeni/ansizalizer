package event

import (
	"github.com/charmbracelet/bubbles/key"
)

type Map struct {
	Enter key.Binding
	Nav   key.Binding
	Right key.Binding
	Left  key.Binding
	Up    key.Binding
	Down  key.Binding
	Copy  key.Binding
	Esc   key.Binding
}

var KeyMap Map

func InitKeyMap() {
	KeyMap = Map{
		Enter: key.NewBinding(
			key.WithKeys("return", "enter"),
			key.WithHelp("↲", "select"),
		),
		Nav: key.NewBinding(
			key.WithKeys("up", "down", "right", "left"),
			key.WithHelp("↕/↔", "navigate"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
		),
		Copy: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "copy to clipboard")),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k Map) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Nav, k.Copy, k.Esc}
}

func (k Map) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Enter, k.Nav, k.Copy, k.Esc}}
}
