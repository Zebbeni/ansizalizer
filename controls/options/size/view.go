package size

import (
	"fmt"

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
	}

	inputStyle = lipgloss.NewStyle().Width(13).AlignHorizontal(lipgloss.Left)

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
		buttons[i] = style.Copy().Width(11).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m Model) drawInputs() string {
	prompt, placeholder := m.getInputColors(WidthForm)
	m.widthInput.PromptStyle = m.widthInput.PromptStyle.Copy().Foreground(prompt)
	m.widthInput.PlaceholderStyle = m.widthInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.widthInput.Focused() == false {
		m.widthInput.Placeholder = fmt.Sprintf("%4d", m.width)
	} else {
		m.widthInput.Placeholder = "    "
	}
	if m.widthInput.Focused() {
		m.widthInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.widthInput.Cursor.SetMode(cursor.CursorHide)
	}

	prompt, placeholder = m.getInputColors(HeightForm)
	m.heightInput.PromptStyle = m.widthInput.PromptStyle.Copy().Foreground(prompt)
	m.heightInput.PlaceholderStyle = m.widthInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.heightInput.Focused() == false {
		m.heightInput.Placeholder = fmt.Sprintf("%4d", m.height)
	} else {
		m.heightInput.Placeholder = "    "
	}
	if m.heightInput.Focused() {
		m.heightInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.heightInput.Cursor.SetMode(cursor.CursorHide)
	}

	width := inputStyle.Render(m.widthInput.View())
	height := inputStyle.Render(m.heightInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, width, height)
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.focus == state {
		return focusColor, focusColor
	} else if m.active == state {
		return activeColor, activeColor
	}
	return normalColor, normalColor
}
