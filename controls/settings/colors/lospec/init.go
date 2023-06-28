package lospec

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
)

var (
	promptStyle      = lipgloss.NewStyle().Padding(0, 1, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

// TODO: This is basically the same as we have in adaptive. Maybe generalize?
func newInput(state State, value string) textinput.Model {
	textinput.New()
	input := textinput.New()
	input.Prompt = stateNames[state]
	input.PromptStyle = promptStyle
	input.PlaceholderStyle = placeholderStyle
	input.Cursor.Blink = true
	input.SetValue(value)
	return input
}

func (m Model) InitializeList() (Model, tea.Cmd) {
	m.didInitializeList = true
	return m.searchLospec(0)
}

func (m Model) DidInitializeList() bool {
	return m.didInitializeList
}
