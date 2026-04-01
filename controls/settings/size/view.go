package size

import (
	"image/color"

	"charm.land/lipgloss/v2"

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
)

func (m Model) drawModeButtons() string {
	title := style.DimmedTitle.PaddingLeft(1).Render("Mode:")

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
	wPrompt, wText, wCursor := m.getInputStyles(WidthForm)
	m.widthInput.SetWidth(3)
	wStyles := m.widthInput.Styles()
	wStyles.Focused.Prompt = wPrompt.Width(8).PaddingLeft(1)
	wStyles.Focused.Text = wText
	wStyles.Blurred.Prompt = wPrompt.Width(8).PaddingLeft(1)
	wStyles.Blurred.Text = wText
	wStyles.Cursor.Blink = m.widthInput.Focused()
	wStyles.Cursor.Color = wCursor
	m.widthInput.SetStyles(wStyles)
	m.widthInput.SetVirtualCursor(m.widthInput.Focused())

	hPrompt, hText, hCursor := m.getInputStyles(HeightForm)
	hStyles := m.heightInput.Styles()
	hStyles.Focused.Prompt = hPrompt.Width(8).PaddingLeft(1)
	hStyles.Focused.Text = hText
	hStyles.Blurred.Prompt = hPrompt.Width(8).PaddingLeft(1)
	hStyles.Blurred.Text = hText
	hStyles.Cursor.Blink = m.heightInput.Focused()
	hStyles.Cursor.Color = hCursor
	m.heightInput.SetStyles(hStyles)
	m.heightInput.SetVirtualCursor(m.heightInput.Focused())

	inputStyle := style.BgStyle().Width(14).AlignHorizontal(lipgloss.Left)
	width := inputStyle.Render(m.widthInput.View())
	height := inputStyle.Render(m.heightInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, width, height)
}

func (m Model) drawCharRatioForm() string {
	crPrompt, crText, crCursor := m.getInputStyles(CharRatioForm)
	m.charRatioInput.SetWidth(30)
	crStyles := m.charRatioInput.Styles()
	crStyles.Focused.Prompt = crPrompt.Width(20).PaddingLeft(1)
	crStyles.Focused.Text = crText
	crStyles.Blurred.Prompt = crPrompt.Width(20).PaddingLeft(1)
	crStyles.Blurred.Text = crText
	crStyles.Cursor.Blink = m.charRatioInput.Focused()
	crStyles.Cursor.Color = crCursor
	m.charRatioInput.SetStyles(crStyles)
	m.charRatioInput.SetVirtualCursor(m.charRatioInput.Focused())

	return style.BgStyle().Width(28).AlignHorizontal(lipgloss.Left).PaddingTop(1).Render(m.charRatioInput.View())
}

func (m Model) getInputStyles(state State) (lipgloss.Style, lipgloss.Style, color.Color) {
	promptStyle := style.DimmedTitle
	textStyle := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if m.IsActive && m.focus == state {
		if m.active == state {
			promptStyle = style.SelectedTitle
			textStyle = lipgloss.NewStyle().Foreground(style.SelectedColor1)
			cursorColor = style.SelectedColor1
		} else {
			promptStyle = style.NormalTitle
			textStyle = lipgloss.NewStyle().Foreground(style.NormalColor1)
			cursorColor = style.NormalColor1
		}
	}
	return promptStyle, textStyle, cursorColor
}
