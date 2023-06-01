package style

import "github.com/charmbracelet/lipgloss"

var (
	NormalTitle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#aaaaaa"})
	NormalParagraph = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#888888"})

	//SelectedTitle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})
	SelectedTitle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#444444", Dark: "#ffffff"})
	SelectedParagraph = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#dddddd"})

	DimmedTitle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#777777"})
	DimmedParagraph = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#aaaaaa", Dark: "#666666"})
)
