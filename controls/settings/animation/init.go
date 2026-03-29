package animation

import (
	"strconv"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(12).Padding(0, 0, 0, 1)
	placeholderStyle = lipgloss.NewStyle()
)

func newInput(label string, value int) textinput.Model {
	input := textinput.New()
	input.Prompt = label
	input.CharLimit = 4
	input.SetValue(strconv.Itoa(value))

	styles := input.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	styles.Focused.Placeholder = placeholderStyle
	styles.Blurred.Placeholder = placeholderStyle
	input.SetStyles(styles)

	return input
}
