package adjust

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(12).Padding(0, 0, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

func newInput(label string, value int) textinput.Model {
	input := textinput.New()
	input.Prompt = label
	input.PromptStyle = promptStyle
	input.PlaceholderStyle = placeholderStyle
	input.CharLimit = 4
	input.SetValue(strconv.Itoa(value))
	return input
}
