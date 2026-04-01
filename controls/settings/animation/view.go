package animation

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inputStyle = style.BgStyle().PaddingLeft(1).Width(21).AlignHorizontal(lipgloss.Center)
)

func (m Model) drawDelayForm() string {
	prompt, text := style.InputColors(m.IsActive, m.focus == DelayForm, m.active == DelayForm)
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
