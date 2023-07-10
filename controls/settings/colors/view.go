package colors

import "github.com/charmbracelet/lipgloss"

var (
	stateOrder = []State{NoPalette, Load, Adapt, Lospec}
	stateNames = map[State]string{
		NoPalette: "True",
		Load:      "Load",
		Adapt:     "Adapt",
		Lospec:    "Lospec",
	}

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
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawTitle() string {
	return titleStyle.Copy().Italic(true).Width(m.width).Align(lipgloss.Center).Render("Colors")
}

func (m Model) drawButtons() string {
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		style := normalStyle
		if m.IsActive && state == m.focus {
			style = focusStyle
		} else if state == m.selected {
			style = activeStyle
		}
		buttons[i] = style.Copy().AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}
