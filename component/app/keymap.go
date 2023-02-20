package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Quit key.Binding
}

func initKeymap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
}

func (a *App) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, a.km.Quit):
		return tea.Quit
	}
	// check with child components here until one returns a non-nil
	// command or all handlers have been hit. Need to figure out
	// how to handle events that affect the app state.
	return nil
}
