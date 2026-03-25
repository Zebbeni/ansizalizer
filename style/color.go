package style

import "github.com/charmbracelet/lipgloss"

// These are all set by Refresh() from the active theme.
// They are kept as exported vars so existing code continues to work.
var (
	NormalColor1   lipgloss.TerminalColor
	NormalColor2   lipgloss.TerminalColor
	SelectedColor1 lipgloss.TerminalColor
	SelectedColor2 lipgloss.TerminalColor
	ExtraDimColor  lipgloss.TerminalColor
	DimmedColor1   lipgloss.TerminalColor
	DimmedColor2   lipgloss.TerminalColor

	NormalTitle     lipgloss.Style
	NormalParagraph lipgloss.Style

	SelectedTitle     lipgloss.Style
	SelectedParagraph lipgloss.Style

	DimmedTitle     lipgloss.Style
	ExtraDimTitle   lipgloss.Style
	DimmedParagraph lipgloss.Style

	ActiveButton lipgloss.Style
	FocusButton  lipgloss.Style
	NormalButton lipgloss.Style

	ActiveButtonNode lipgloss.Style
	FocusButtonNode  lipgloss.Style
	NormalButtonNode lipgloss.Style

	InactiveTabBorder lipgloss.Border
	ActiveTabBorder   lipgloss.Border
	InactiveTabStyle  lipgloss.Style
	ActiveTabStyle    lipgloss.Style
	TabWindowStyle    lipgloss.Style
)

func init() {
	InactiveTabBorder = TabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder = TabBorderWithBottom("┘", " ", "└")
	Refresh()
}

func TabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
