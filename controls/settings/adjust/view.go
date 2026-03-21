package adjust

import (
	"strings"

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
	m.brightnessInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(text)
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
	m.contrastInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(text)
	if m.contrastInput.Focused() {
		m.contrastInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.contrastInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Render(m.contrastInput.View())
}

func (m Model) drawSlider(state State, value int) string {
	sliderWidth := m.width - 2
	if sliderWidth < 3 {
		sliderWidth = 3
	}

	// Map value from [-100, 100] to position [0, sliderWidth-1]
	pos := int(float64(value+100) / 200.0 * float64(sliderWidth-1))
	if pos < 0 {
		pos = 0
	}
	if pos >= sliderWidth {
		pos = sliderWidth - 1
	}

	left := strings.Repeat("-", pos)
	right := strings.Repeat("-", sliderWidth-1-pos)

	sliderColor := normalColor
	if m.IsActive && m.focus == state {
		sliderColor = focusColor
	}

	style := lipgloss.NewStyle().Foreground(sliderColor).PaddingLeft(1)
	return style.Render(left + "|" + right)
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
