package adaptive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(8).PaddingLeft(1)
	placeholderStyle = lipgloss.NewStyle()
)

func newInput(state State) textinput.Model {
	textinput.New()
	input := textinput.New()
	input.Prompt = stateNames[state]
	input.PromptStyle = promptStyle
	input.PlaceholderStyle = placeholderStyle
	input.Cursor.Blink = true
	input.CharLimit = 3
	input.SetValue(fmt.Sprintf("16"))
	return input
}
