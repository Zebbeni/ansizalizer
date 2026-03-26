package animation

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inputStyle = style.BgStyle().Width(20).AlignHorizontal(lipgloss.Left)
)

func (m Model) drawDelayForm() string {
	prompt, text := m.getInputColors(DelayForm)
	m.delayInput.Width = 5
	m.delayInput.PromptStyle = m.delayInput.PromptStyle.Copy().Foreground(prompt)
	m.delayInput.TextStyle = m.delayInput.TextStyle.Copy().Foreground(text)
	m.delayInput.Cursor.TextStyle = style.BgStyle().Foreground(text)
	if m.delayInput.Focused() {
		m.delayInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.delayInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Render(m.delayInput.View())
}

func (m Model) getInputColors(state State) (lipgloss.TerminalColor, lipgloss.TerminalColor) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return style.NormalColor1, style.SelectedColor1
		}
		return style.SelectedColor1, style.NormalColor1
	}
	return style.DimmedColor1, style.DimmedColor1
}
