package adaptive

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder = []State{CountForm, IterForm}
	stateNames = map[State]string{
		CountForm: "Colors",
		IterForm:  "Passes",
	}

	inputStyle = lipgloss.NewStyle().Width(13).AlignHorizontal(lipgloss.Left)

	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
	titleStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawTitle() string {
	title := style.DimmedTitle.Copy().Italic(true).Render("Create palette From image")
	return lipgloss.NewStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}

func (m Model) drawInputs() string {
	prompt, placeholder := m.getInputColors(CountForm)

	m.countInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.countInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.countInput.Focused() {
		m.countInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.countInput.Cursor.SetMode(cursor.CursorHide)
	}

	prompt, placeholder = m.getInputColors(IterForm)
	m.iterInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.iterInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.iterInput.Focused() {
		m.iterInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.iterInput.Cursor.SetMode(cursor.CursorHide)
	}

	countInput := inputStyle.Render(m.countInput.View())
	iterInput := inputStyle.Render(m.iterInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, countInput, iterInput)
}

func (m Model) drawGenerateButton() string {
	styleColor := normalColor
	if m.IsActive && m.focus == Generate {
		styleColor = focusColor
	} else if m.active == Generate {
		styleColor = activeColor
	}

	style := lipgloss.NewStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styleColor).
		Foreground(styleColor)

	button := style.Render("Generate New")
	return lipgloss.NewStyle().Width(m.width - 2).AlignHorizontal(lipgloss.Center).Render(button)
}

// TODO: This is almost the same as drawGenerateButton. See if we can generalize
func (m Model) drawSaveButton() string {
	styleColor := normalColor
	if m.IsActive && m.focus == Save {
		styleColor = focusColor
	} else if m.active == Save {
		styleColor = activeColor
	}

	style := lipgloss.NewStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		PaddingTop(1).
		Foreground(styleColor)

	button := style.Render("Save to .hex File")
	return lipgloss.NewStyle().Width(m.width - 2).AlignHorizontal(lipgloss.Center).Render(button)
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.IsActive {
		if m.focus == state {
			return focusColor, focusColor
		} else if m.active == state {
			return activeColor, activeColor
		}
	}
	return normalColor, normalColor
}
