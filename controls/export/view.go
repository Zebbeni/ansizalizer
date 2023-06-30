package export

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
)

func (m Model) renderWithBorder(content string, state State) string {
	renderColor := normalColor
	if m.active == state {
		renderColor = activeColor
	} else if m.focus == state {
		renderColor = focusColor
	}

	textStyle := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1, 0, 1).
		Foreground(renderColor)
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(renderColor)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: textStyle,
	}

	return renderer.Render(stateTitles[state], content, m.width-2)
}
