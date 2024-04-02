package size

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"
)

var (
	stateOrder = []State{FitButton, StretchButton}
	stateNames = map[State]string{
		FitButton:     "Fit",
		StretchButton: "Stretch",
		WidthForm:     "Width",
		HeightForm:    "Height",
		CharRatioForm: "Char Size Ratio (Width/Height)",
	}

	inputStyle = lipgloss.NewStyle().Width(14).AlignHorizontal(lipgloss.Left)

	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
	titleStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawButtons() string {
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		styleColor := normalColor
		if m.IsActive {
			if state == m.focus {
				styleColor = focusColor
			} else if state == m.active {
				styleColor = activeColor
			}
		}
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styleColor).
			Foreground(styleColor)
		buttons[i] = style.Copy().Width(12).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m Model) drawSizeForms() string {
	prompt, text := m.getInputColors(WidthForm)
	m.widthInput.Width = 3
	m.widthInput.PromptStyle = m.widthInput.PromptStyle.Copy().Foreground(prompt)
	m.heightInput.TextStyle = m.heightInput.TextStyle.Copy().Foreground(text)
	if m.widthInput.Focused() {
		m.widthInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.widthInput.Cursor.SetMode(cursor.CursorHide)
	}

	prompt, text = m.getInputColors(HeightForm)
	m.heightInput.PromptStyle = m.heightInput.PromptStyle.Copy().Foreground(prompt)
	m.heightInput.TextStyle = m.heightInput.TextStyle.Copy().Foreground(text)
	if m.heightInput.Focused() {
		m.heightInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.heightInput.Cursor.SetMode(cursor.CursorHide)
	}

	width := inputStyle.Render(m.widthInput.View())
	height := inputStyle.Render(m.heightInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, width, height)
}

func (m Model) drawCharRatioForm() string {
	prompt, text := m.getInputColors(CharRatioForm)
	m.charRatioInput.Width = 30
	m.charRatioInput.PromptStyle = m.charRatioInput.PromptStyle.Copy().Width(20).Foreground(prompt)
	m.charRatioInput.TextStyle = m.charRatioInput.TextStyle.Copy().Foreground(text)
	if m.charRatioInput.Focused() {
		m.charRatioInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.charRatioInput.Cursor.SetMode(cursor.CursorHide)
	}

	return inputStyle.Copy().Width(28).AlignHorizontal(lipgloss.Left).PaddingTop(1).Render(m.charRatioInput.View())
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.focus == state {
		if m.active == state {
			return activeColor, focusColor
		} else {
			return focusColor, activeColor
		}
	}
	return normalColor, normalColor
}
