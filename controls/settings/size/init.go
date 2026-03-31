package size

import (
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(8).Padding(0, 0, 0, 1)
	placeholderStyle = lipgloss.NewStyle()

	floatPromptStyle      = lipgloss.NewStyle().Padding(0, 1)
	floatPlaceholderStyle = lipgloss.NewStyle()
)


func newInput(state State, value int) textinput.Model {
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
	input.CharLimit = 3
	input.SetValue(strconv.Itoa(value))
	return input
}

func newFloatInput(state State, value float64) textinput.Model {
	input := textinput.New()
	input.Prompt = stateNames[state]
	styles := input.Styles()
	styles.Focused.Prompt = floatPromptStyle
	styles.Blurred.Prompt = floatPromptStyle
	styles.Focused.Placeholder = floatPlaceholderStyle
	styles.Blurred.Placeholder = floatPlaceholderStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	input.SetStyles(styles)
	input.CharLimit = 5
	input.SetValue(fmt.Sprintf("%1.2f", value))
	return input
}
