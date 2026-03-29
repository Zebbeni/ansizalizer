package export

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) renderWithBorder(content string, state State) string {
	renderColor := style.DimmedColor1
	if m.active == state {
		renderColor = style.NormalColor1
	} else if m.focus == state {
		renderColor = style.SelectedColor1
	}

	textStyle := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1, 0, 1).
		Foreground(renderColor)
	border := lipgloss.RoundedBorder()
	if m.focus == state && m.active != state {
		border = style.HeavyBorder()
	}
	borderStyle := style.BgStyle().
		Border(border).
		BorderForeground(renderColor).
		BorderBackground(style.ActiveTheme.Bg)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: textStyle,
	}

	return renderer.Render(stateTitles[state], content, m.width-2)
}

func (m Model) drawProcessButton() string {
	buttonStyle := style.NormalButton
	if m.focus == Process {
		buttonStyle = style.FocusButton
	}

	centerStyle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)

	internalStyle := centerStyle.Copy().Width(m.width - 2)
	title := internalStyle.Render(stateTitles[Process])
	button := buttonStyle.Render(title)

	return centerStyle.Copy().Width(m.width).AlignHorizontal(lipgloss.Center).Render(button)
}
