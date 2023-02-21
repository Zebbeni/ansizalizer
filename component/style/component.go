package style

import "github.com/charmbracelet/lipgloss"

var (
	Border = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	Text   = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
)
