package characters

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/style"
)

func newInput(prompt string, value string) textinput.Model {
	input := textinput.New()
	input.Prompt = prompt

	styles := input.Styles()
	styles.Focused.Prompt = style.NormalButtonNode.Padding(0, 1, 0, 0)
	styles.Blurred.Prompt = style.NormalButtonNode.Padding(0, 1, 0, 0)
	styles.Focused.Placeholder = lipgloss.NewStyle()
	styles.Blurred.Placeholder = lipgloss.NewStyle()
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	input.SetStyles(styles)

	input.SetValue(value)
	return input
}

func newNumberInput(prompt string, isFloat bool, def float64, min *float64) numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:  prompt,
		IsFloat: isFloat,
		Min:     min,
		Default: def,
	})
	styles := m.Styles()
	styles.Focused.Prompt = style.NormalButtonNode.Padding(0, 1, 0, 0)
	styles.Blurred.Prompt = style.NormalButtonNode.Padding(0, 1, 0, 0)
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	m.SetStyles(styles)
	return m
}
