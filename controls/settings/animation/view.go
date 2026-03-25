package animation

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")

	inputStyle = lipgloss.NewStyle().Width(20).AlignHorizontal(lipgloss.Left)
)

func (m Model) drawDelayForm() string {
	prompt, text := m.getInputColors(DelayForm)
	m.delayInput.Width = 5
	m.delayInput.PromptStyle = m.delayInput.PromptStyle.Copy().Foreground(prompt)
	m.delayInput.TextStyle = m.delayInput.TextStyle.Copy().Foreground(text)
	m.delayInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(text)
	if m.delayInput.Focused() {
		m.delayInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.delayInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Render(m.delayInput.View())
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return activeColor, focusColor
		}
		return focusColor, activeColor
	}
	return normalColor, normalColor
}
