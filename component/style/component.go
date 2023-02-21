package style

import "github.com/charmbracelet/lipgloss"

var (
	ViewportBorder = lipgloss.NewStyle()
	ControlsBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	ViewerBorder   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	Text           = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	Help           = lipgloss.NewStyle().Margin(0, 0, 0, 1)
)
