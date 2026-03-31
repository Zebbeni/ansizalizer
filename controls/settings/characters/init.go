package characters

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

// TODO: This is basically the same as we have in adaptive. Maybe generalize?
func newInput(prompt string, value string) textinput.Model {
	input := textinput.New()
	input.Prompt = prompt

	styles := input.Styles()
	styles.Focused.Prompt = style.NormalButtonNode.Copy().Padding(0, 1, 0, 0)
	styles.Blurred.Prompt = style.NormalButtonNode.Copy().Padding(0, 1, 0, 0)
	styles.Focused.Placeholder = lipgloss.NewStyle()
	styles.Blurred.Placeholder = lipgloss.NewStyle()
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	input.SetStyles(styles)

	input.SetValue(value)
	return input
}
