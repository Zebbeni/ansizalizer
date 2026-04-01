package lospec

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"charm.land/bubbles/v2/textinput"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/style"
)

var (
	promptStyle      = lipgloss.NewStyle().Padding(0, 1, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

func newCountInput() numberinput.Model {
	return numberinput.New(numberinput.Options{
		Prompt:    stateNames[CountForm] + " ",
		CharLimit: 3,
		IsFloat:   false,
		Min:       numberinput.FloatPtr(1),
		Default:   16,
	})
}

func newInput(state State, value string) textinput.Model {
	input := textinput.New()
	input.Prompt = stateNames[state]
	styles := input.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	styles.Focused.Placeholder = placeholderStyle
	styles.Blurred.Placeholder = placeholderStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
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
