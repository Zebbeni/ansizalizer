package characters

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle      = lipgloss.NewStyle().Padding(0, 1, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

// TODO: This is basically the same as we have in adaptive. Maybe generalize?
func newInput(prompt string, value string) textinput.Model {
	textinput.New()
	input := textinput.New()
	input.Prompt = prompt
	input.PromptStyle = promptStyle
	input.PlaceholderStyle = placeholderStyle
	input.Cursor.Blink = true
	input.SetValue(value)
	return input
}
