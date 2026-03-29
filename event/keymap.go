package event

import (
	"charm.land/bubbles/v2/key"
)

type Map struct {
	Enter    key.Binding
	Nav      key.Binding
	Right    key.Binding
	Left     key.Binding
	Up       key.Binding
	Down     key.Binding
	Expand   key.Binding
	Collapse key.Binding
	Copy     key.Binding
	Save     key.Binding
	Debug    key.Binding
	Esc      key.Binding
}

var KeyMap Map

func InitKeyMap() {
	KeyMap = Map{
		Enter: key.NewBinding(
			key.WithKeys("return", "enter"),
			key.WithHelp("↲/enter", "select/expand menu"),
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
		Expand: key.NewBinding(
			key.WithKeys("+", "="),
			key.WithHelp("+", "expand menu"),
		),
		Collapse: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "collapse menu"),
		),
		Copy: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "copy to clipboard")),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save to file")),
		Debug: key.NewBinding(
			key.WithKeys("ctrl+d")),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/exit menu"),
		),
	}
}

func (k Map) ShortHelp() []key.Binding {
	return []key.Binding{k.Nav, k.Enter, k.Esc, k.Expand, k.Collapse, k.Copy, k.Save}
}

func (k Map) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Nav, k.Enter, k.Esc, k.Expand, k.Collapse, k.Copy, k.Save}}
}
