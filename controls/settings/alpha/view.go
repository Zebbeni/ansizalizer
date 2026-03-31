package alpha

import (
	"charm.land/lipgloss/v2"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawAlphaOptions() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Enable:")

	yesButton := style.NormalButtonNode
	if m.IsActive && m.focus == AlphaYes {
		yesButton = style.FocusButtonNode
	} else if m.useAlpha == true {
		yesButton = style.ActiveButtonNode
	}
	yesNode := yesButton.Render("Yes")
	yesNode = style.BgStyle().PaddingLeft(1).Render(yesNode)

	noButton := style.NormalButtonNode
	if m.IsActive && m.focus == AlphaNo {
		noButton = style.FocusButtonNode
	} else if m.useAlpha == false {
		noButton = style.ActiveButtonNode
	}
	noNode := noButton.Render("No")
	noNode = style.BgStyle().PaddingLeft(1).Render(noNode)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, yesNode, noNode)
}

func (m Model) drawThresholdInput() string {
	promptStyle := style.DimmedTitle.Copy()
	textSt := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if m.thresholdInput.Focused() {
		promptStyle = style.SelectedTitle.Copy()
		textSt = lipgloss.NewStyle().Foreground(style.SelectedColor1)
		cursorColor = style.SelectedColor1
	} else if m.focus == ThresholdForm && m.IsActive {
		promptStyle = style.NormalTitle.Copy()
		textSt = lipgloss.NewStyle().Foreground(style.NormalColor1)
		cursorColor = style.NormalColor1
	}
	styles := m.thresholdInput.Styles()
	styles.Focused.Prompt = promptStyle.PaddingLeft(1)
	styles.Focused.Text = textSt
	styles.Blurred.Prompt = promptStyle.PaddingLeft(1)
	styles.Blurred.Text = textSt
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.thresholdInput.SetStyles(styles)
	m.thresholdInput.SetVirtualCursor(m.thresholdInput.Focused())
	return m.thresholdInput.View()
}

func (m Model) drawAlphaTrimOptions() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Trim Output:")

	yesButton := style.NormalButtonNode
	if m.IsActive && m.focus == TrimAlphaYes {
		yesButton = style.FocusButtonNode
	} else if m.trimAlpha == true {
		yesButton = style.ActiveButtonNode
	}
	yesNode := yesButton.Render("Yes")
	yesNode = style.BgStyle().PaddingLeft(1).Render(yesNode)

	noButton := style.NormalButtonNode
	if m.IsActive && m.focus == TrimAlphaNo {
		noButton = style.FocusButtonNode
	} else if m.trimAlpha == false {
		noButton = style.ActiveButtonNode
	}
	noNode := noButton.Render("No")
	noNode = style.BgStyle().PaddingLeft(1).Render(noNode)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, yesNode, noNode)
}
