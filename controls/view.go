package controls

import "github.com/charmbracelet/lipgloss"

var (
	activeStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#aaaaaa")).
			Foreground(lipgloss.Color("#aaaaaa"))
	focusStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#ffffff"))
	normalStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555")).
			Foreground(lipgloss.Color("#555555"))
)

// ▛▜▐▜▐▀▐▝▜▐▜▐ ▐▝▜▐▛▐▀▌
// ▛▜▐▐▗▟▐▐▄▐▜▐▄▐▐▄▐▄▐▜
func (m Model) drawTitle() string {
	title1Runes := []rune{' ', '▛', '▜', '▐', '▜', '▐', '▀', '▐', '▝', '▜', '▐', '▜', '▐', ' ', '▐', '▝', '▜', '▐', '▛', '▐', '▀', '▌'}
	title2Runes := []rune{'▛', '▜', '▐', '▐', '▗', '▟', '▐', '▐', '▄', '▐', '▜', '▐', '▄', '▐', '▐', '▄', '▐', '▄', '▐', '▜'}
	title1 := lipgloss.JoinHorizontal(lipgloss.Left, string(title1Runes))
	title2 := lipgloss.JoinHorizontal(lipgloss.Left, string(title2Runes))
	title := lipgloss.JoinVertical(lipgloss.Left, title1, title2)
	return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Padding(1, 0, 0, 0).Render(title)
}

func (m Model) drawButtons() string {
	buttonWidth := (m.width - (numButtons * 2)) / numButtons
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		style := normalStyle
		if state == m.active {
			style = activeStyle
		} else if state == m.focus {
			style = focusStyle
		}
		buttons[i] = style.Copy().Width(buttonWidth).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	return buttonRow
}
