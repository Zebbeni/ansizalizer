package animation

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inputStyle = style.BgStyle().Width(20).AlignHorizontal(lipgloss.Left)
)

func (m Model) drawDelayForm() string {
	prompt, text := m.getInputColors(DelayForm)
	m.delayInput.SetWidth(5)

	styles := m.delayInput.Styles()
	styles.Focused.Prompt = styles.Focused.Prompt.Foreground(prompt)
	styles.Blurred.Prompt = styles.Blurred.Prompt.Foreground(prompt)
	styles.Focused.Text = styles.Focused.Text.Foreground(text)
	styles.Blurred.Text = styles.Blurred.Text.Foreground(text)
	styles.Cursor.Color = text
	styles.Cursor.Blink = true
	m.delayInput.SetStyles(styles)
	m.delayInput.SetVirtualCursor(m.delayInput.Focused())

	return inputStyle.Render(m.delayInput.View())
}

func (m Model) getInputColors(state State) (color.Color, color.Color) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return style.NormalColor1, style.SelectedColor1
		}
		return style.SelectedColor1, style.NormalColor1
	}
	return style.DimmedColor1, style.DimmedColor1
}
