package lospec

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/textinput"
)

var (
	promptStyle      = lipgloss.NewStyle().Padding(0, 1, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

// TODO: This is basically the same as we have in adaptive. Maybe generalize?
func newInput(state State, value string) textinput.Model {
	input := textinput.New()
	input.Prompt = stateNames[state]
	styles := input.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	styles.Focused.Placeholder = placeholderStyle
	styles.Blurred.Placeholder = placeholderStyle
	styles.Cursor.Blink = true
	input.SetStyles(styles)
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
