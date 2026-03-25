package size

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder = []State{FitButton, FillButton, StretchButton}
	stateNames = map[State]string{
		FitButton:     "Fit",
		FillButton:    "Fill",
		StretchButton: "Stretch",
		WidthForm:     "Width",
		HeightForm:    "Height",
		CharRatioForm: "Char Size Ratio (Width/Height)",
	}

	inputStyle = style.BgStyle().Width(14).AlignHorizontal(lipgloss.Left)
)

func (m Model) drawModeButtons() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Mode:")

	fitStyle := style.NormalButtonNode
	if m.IsActive && FitButton == m.focus {
		fitStyle = style.FocusButtonNode
	} else if m.mode == Fit {
		fitStyle = style.ActiveButtonNode
	}
	fitButton := fitStyle.Render("Fit")
	fitButton = style.BgStyle().Width(6).AlignHorizontal(lipgloss.Center).Render(fitButton)

	fillStyle := style.NormalButtonNode
	if m.IsActive && FillButton == m.focus {
		fillStyle = style.FocusButtonNode
	} else if m.mode == Fill {
		fillStyle = style.ActiveButtonNode
	}
	fillButton := fillStyle.Render("Fill")
	fillButton = style.BgStyle().Width(6).AlignHorizontal(lipgloss.Center).Render(fillButton)

	stretchStyle := style.NormalButtonNode
	if m.IsActive && StretchButton == m.focus {
		stretchStyle = style.FocusButtonNode
	} else if m.mode == Stretch {
		stretchStyle = style.ActiveButtonNode
	}
	stretchButton := stretchStyle.Render("Stretch")
	stretchButton = style.BgStyle().Width(9).AlignHorizontal(lipgloss.Center).Render(stretchButton)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, fitButton, fillButton, stretchButton)
}

func (m Model) drawSizeForms() string {
	prompt, text := m.getInputColors(WidthForm)
	m.widthInput.Width = 3
	m.widthInput.PromptStyle = m.widthInput.PromptStyle.Copy().Foreground(prompt)
	m.widthInput.TextStyle = m.widthInput.TextStyle.Copy().Foreground(text)
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

func (m Model) getInputColors(state State) (lipgloss.TerminalColor, lipgloss.TerminalColor) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return style.NormalColor1, style.SelectedColor1
		}
		return style.SelectedColor1, style.NormalColor1
	}
	return style.DimmedColor1, style.DimmedColor1
}
