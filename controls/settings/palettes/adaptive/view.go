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
)

func (m Model) drawTitle() string {
	title := style.DimmedTitle.Copy().Italic(true).Render("Generate palette From image")
	return style.BgStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}

func (m Model) drawInputs() string {
	promptStyle, textStyle := m.getInputStyles(CountForm)
	m.countInput.PromptStyle = promptStyle
	m.countInput.TextStyle = textStyle
	if m.countInput.Focused() {
		m.countInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.countInput.Cursor.SetMode(cursor.CursorHide)
	}

	promptStyle, textStyle = m.getInputStyles(IterForm)
	m.iterInput.PromptStyle = promptStyle
	m.iterInput.TextStyle = textStyle
	if m.iterInput.Focused() {
		m.iterInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.iterInput.Cursor.SetMode(cursor.CursorHide)
	}

	inputStyle := style.BgStyle().Width(13).AlignHorizontal(lipgloss.Left).PaddingLeft(1)
	countInput := inputStyle.Render(m.countInput.View())
	iterInput := inputStyle.Render(m.iterInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, countInput, iterInput)
}

func (m Model) drawGenerateButton() string {
	styleColor := style.DimmedColor1
	if m.IsActive && m.focus == Generate {
		styleColor = style.SelectedColor1
	} else if m.active == Generate {
		styleColor = style.NormalColor1
	}

	btnStyle := style.BgStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styleColor).
		Foreground(styleColor)

	button := btnStyle.Render("Generate New")
	return style.BgStyle().Width(m.width).PaddingLeft(1).Render(button)
}

// TODO: This is almost the same as drawGenerateButton. See if we can generalize
func (m Model) drawSaveButton() string {
	styleColor := style.DimmedColor1
	if m.IsActive && m.focus == Save {
		styleColor = style.SelectedColor1
	} else if m.active == Save {
		styleColor = style.NormalColor1
	}

	btnStyle := style.BgStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		PaddingTop(1).
		Foreground(styleColor)

	button := btnStyle.Render("Save to .hex File")
	return style.BgStyle().Width(m.width).PaddingLeft(1).Render(button)
}

func (m Model) getInputStyles(state State) (lipgloss.Style, lipgloss.Style) {
	promptStyle := style.DimmedTitle.Copy()
	textStyle := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	if m.IsActive && m.focus == state {
		promptStyle = style.SelectedTitle.Copy()
		textStyle = lipgloss.NewStyle().Foreground(style.SelectedColor1)
	} else if m.active == state {
		promptStyle = style.NormalTitle.Copy()
		textStyle = lipgloss.NewStyle().Foreground(style.NormalColor1)
	}
	return promptStyle.Width(8).PaddingLeft(1), textStyle
}
