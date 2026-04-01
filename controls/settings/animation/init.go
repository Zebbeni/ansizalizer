package animation

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
)

var (
	promptStyle = lipgloss.NewStyle().Width(12).Padding(0, 0, 0, 1)
)

func newInput(label string, value int) numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:    label + " ",
		CharLimit: 4,
		IsFloat:   false,
		Min:       numberinput.FloatPtr(10),
		Max:       numberinput.FloatPtr(2000),
		Default:   float64(value),
	})
	styles := m.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	m.SetStyles(styles)
	return m
}
