package adaptive

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateNames = map[State]string{
		CountForm: "Colors",
		IterForm:  "Passes",
	}
)

func (m Model) drawInputs() string {
	promptStyle, textStyle := m.getInputStyles(CountForm)
	countStyles := m.countInput.Styles()
	countStyles.Focused.Prompt = promptStyle
	countStyles.Focused.Text = textStyle
	countStyles.Blurred.Prompt = promptStyle
	countStyles.Blurred.Text = textStyle
	countStyles.Cursor.Blink = m.countInput.Focused()
	m.countInput.SetStyles(countStyles)
	m.countInput.SetVirtualCursor(m.countInput.Focused())

	promptStyle, textStyle = m.getInputStyles(IterForm)
	iterStyles := m.iterInput.Styles()
	iterStyles.Focused.Prompt = promptStyle
	iterStyles.Focused.Text = textStyle
	iterStyles.Blurred.Prompt = promptStyle
	iterStyles.Blurred.Text = textStyle
	iterStyles.Cursor.Blink = m.iterInput.Focused()
	m.iterInput.SetStyles(iterStyles)
	m.iterInput.SetVirtualCursor(m.iterInput.Focused())

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
	promptStyle := style.DimmedTitle
	textStyle := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	if m.IsActive && m.focus == state {
		promptStyle = style.SelectedTitle
		textStyle = lipgloss.NewStyle().Foreground(style.SelectedColor1)
	} else if m.active == state {
		promptStyle = style.NormalTitle
		textStyle = lipgloss.NewStyle().Foreground(style.NormalColor1)
	}
	return promptStyle.Width(8).PaddingLeft(1), textStyle
}
