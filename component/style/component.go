package style

import "github.com/charmbracelet/lipgloss"

var (
	ViewportBorder = lipgloss.NewStyle()
	ControlsBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1, 0, 1)
	ViewerBorder   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1, 0, 1)
	Text           = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
	Help           = lipgloss.NewStyle().Inline(true).Margin(0, 0, 0, 1)
	InactiveItem   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	ActiveItem     = lipgloss.NewStyle().Foreground(lipgloss.Color("64")).Background(lipgloss.Color("20"))
)
