package size

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/textinput"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(8).Padding(0, 0, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

func newInput(state State, value int) textinput.Model {
	textinput.New()
	input := textinput.New()
	input.Prompt = stateNames[state]
	input.PromptStyle = promptStyle
	input.PlaceholderStyle = placeholderStyle
	input.CharLimit = 3
	input.SetValue(strconv.Itoa(value))
	return input
}
