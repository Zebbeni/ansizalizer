package settings

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) renderWithBorder(content string, state State, collapsed bool) string {
	renderColor := style.DimmedColor1
	if m.IsActive && m.active == state {
		renderColor = style.NormalColor1
	} else if m.IsActive && m.focus == state {
		renderColor = style.SelectedColor1
	}

	textStyle := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1, 0, 1).
		Foreground(renderColor)
	border := lipgloss.RoundedBorder()
	if m.IsActive && m.focus == state && m.active != state {
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

	title := stateTitles[state]
	if collapsed {
		title += " (+)"
	} else {
		title += " (-)"
	}

	return renderer.Render(title, content, m.width-2)
}
