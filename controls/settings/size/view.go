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
	m.widthInput.SetWidth(3)
	wStyles := m.widthInput.Styles()
	wStyles.Focused.Prompt = wStyles.Focused.Prompt.Foreground(prompt)
	wStyles.Focused.Text = wStyles.Focused.Text.Foreground(text)
	wStyles.Blurred.Prompt = wStyles.Blurred.Prompt.Foreground(prompt)
	wStyles.Blurred.Text = wStyles.Blurred.Text.Foreground(text)
	wStyles.Cursor.Blink = m.widthInput.Focused()
	m.widthInput.SetStyles(wStyles)

	prompt, text = m.getInputColors(HeightForm)
	hStyles := m.heightInput.Styles()
	hStyles.Focused.Prompt = hStyles.Focused.Prompt.Foreground(prompt)
	hStyles.Focused.Text = hStyles.Focused.Text.Foreground(text)
	hStyles.Blurred.Prompt = hStyles.Blurred.Prompt.Foreground(prompt)
	hStyles.Blurred.Text = hStyles.Blurred.Text.Foreground(text)
	hStyles.Cursor.Blink = m.heightInput.Focused()
	m.heightInput.SetStyles(hStyles)

	width := inputStyle.Render(m.widthInput.View())
	height := inputStyle.Render(m.heightInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, width, height)
}

func (m Model) drawCharRatioForm() string {
	prompt, text := m.getInputColors(CharRatioForm)
	m.charRatioInput.SetWidth(30)
	crStyles := m.charRatioInput.Styles()
	crStyles.Focused.Prompt = crStyles.Focused.Prompt.Width(20).Foreground(prompt)
	crStyles.Focused.Text = crStyles.Focused.Text.Foreground(text)
	crStyles.Blurred.Prompt = crStyles.Blurred.Prompt.Width(20).Foreground(prompt)
	crStyles.Blurred.Text = crStyles.Blurred.Text.Foreground(text)
	crStyles.Cursor.Blink = m.charRatioInput.Focused()
	m.charRatioInput.SetStyles(crStyles)

	return inputStyle.Copy().Width(28).AlignHorizontal(lipgloss.Left).PaddingTop(1).Render(m.charRatioInput.View())
}

func (m Model) getInputColors(state State) (color.Color, color.Color) {
	if m.IsActive && m.focus == state {
		if m.active == state {
			return style.NormalColor1, style.SelectedColor1
		}
		return style.SelectedColor1, style.NormalColor1
	}
	return style.DimmedColor1, style.DimmedColor1
}
