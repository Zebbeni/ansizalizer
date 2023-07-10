package controls

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

// ▛▜▐▜▐▀▐▝▜▐▜▐ ▐▝▜▐▛▐▀▌
// ▛▜▐▐▗▟▐▐▄▐▜▐▄▐▐▄▐▄▐▜
func (m Model) drawTitle() string {
	title1Runes := []rune{' ', ' ', '▛', '▜', '▐', '▜', '▐', '▀', '▐', '▝', '▜', '▐', '▜', '▐', ' ', '▐', '▝', '▜', '▐', '▛', '▐', '▀', '▌', ' '}
	title2Runes := []rune{' ', '▛', '▜', '▐', '▐', '▗', '▟', '▐', '▐', '▄', '▐', '▜', '▐', '▄', '▐', '▐', '▄', '▐', '▄', '▐', '▀', '▖'}
	title1 := lipgloss.JoinHorizontal(lipgloss.Left, string(title1Runes))
	title2 := lipgloss.JoinHorizontal(lipgloss.Left, string(title2Runes))
	title := lipgloss.JoinVertical(lipgloss.Left, title1, title2)
	return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Padding(1, 0, 0, 0).Render(title)
}

func (m Model) drawButtons() string {
	buttonWidth := (m.width - (numButtons * 2)) / numButtons
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		buttonStyle := style.NormalButton
		if state == m.active {
			buttonStyle = style.ActiveButton
		} else if state == m.focus {
			buttonStyle = style.FocusButton
		}
		buttons[i] = buttonStyle.Copy().Width(buttonWidth).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	return buttonRow
}

func (m Model) drawBrowserTitle() string {
	title := style.DimmedTitle.Copy().Italic(true).Render("Search Images")
	return lipgloss.NewStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}
