package size

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/style"
)

var (
	promptStyle      = lipgloss.NewStyle().Width(8).Padding(0, 0, 0, 1)
	floatPromptStyle = lipgloss.NewStyle().Padding(0, 1)
)

func newInput(state State, value int) numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:    stateNames[state] + " ",
		CharLimit: 3,
		IsFloat:   false,
		Min:       numberinput.FloatPtr(1),
		Default:   float64(value),
	})
	styles := m.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	m.SetStyles(styles)
	return m
}

func newFloatInput(state State, value float64) numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:    stateNames[state] + " ",
		CharLimit: 5,
		IsFloat:   true,
		Min:       numberinput.FloatPtr(0.01),
		Default:   value,
	})
	styles := m.Styles()
	styles.Focused.Prompt = floatPromptStyle
	styles.Blurred.Prompt = floatPromptStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	m.SetStyles(styles)
	return m
}
