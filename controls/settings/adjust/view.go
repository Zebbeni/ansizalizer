package adjust

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

func (m Model) drawBrightnessForm() string {
	prompt, text := m.getInputColors(BrightnessForm)
	m.brightnessInput.Width = 4
	m.brightnessInput.PromptStyle = m.brightnessInput.PromptStyle.Copy().Foreground(prompt)
	m.brightnessInput.TextStyle = m.brightnessInput.TextStyle.Copy().Foreground(text)
	if m.brightnessInput.Focused() {
		m.brightnessInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.brightnessInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Render(m.brightnessInput.View())
}

func (m Model) drawContrastForm() string {
	prompt, text := m.getInputColors(ContrastForm)
	m.contrastInput.Width = 4
	m.contrastInput.PromptStyle = m.contrastInput.PromptStyle.Copy().Foreground(prompt)
	m.contrastInput.TextStyle = m.contrastInput.TextStyle.Copy().Foreground(text)
	if m.contrastInput.Focused() {
		m.contrastInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.contrastInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Render(m.contrastInput.View())
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.focus == state {
		if m.active == state {
			return activeColor, focusColor
		}
		return focusColor, activeColor
	}
	return normalColor, normalColor
}
