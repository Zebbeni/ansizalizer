package alpha

import (
	"charm.land/lipgloss/v2"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawAlphaOptions() string {
	focused := m.IsActive && m.focus == UseAlpha
	checkbox := style.RenderCheckbox("Make Transparent", m.useAlpha, focused)
	return lipgloss.NewStyle().PaddingLeft(1).Render(checkbox)
}

func (m Model) drawThresholdInput() string {
	promptStyle := style.DimmedTitle
	textSt := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if m.thresholdInput.Focused() {
		promptStyle = style.SelectedTitle
		textSt = lipgloss.NewStyle().Foreground(style.SelectedColor1)
		cursorColor = style.SelectedColor1
	} else if m.focus == ThresholdForm && m.IsActive {
		promptStyle = style.NormalTitle
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
	focused := m.IsActive && m.focus == TrimAlpha
	checkbox := style.RenderCheckbox("Trim Output", m.trimAlpha, focused)
	return lipgloss.NewStyle().PaddingLeft(1).Render(checkbox)
}
