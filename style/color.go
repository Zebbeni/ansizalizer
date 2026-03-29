package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// These are all set by Refresh() from the active theme.
// They are kept as exported vars so existing code continues to work.
var (
	NormalColor1   color.Color
	NormalColor2   color.Color
	SelectedColor1 color.Color
	SelectedColor2 color.Color
	ExtraDimColor  color.Color
	DimmedColor1   color.Color
	DimmedColor2   color.Color

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

	InactiveTabBorder        lipgloss.Border
	ActiveTabBorder          lipgloss.Border
	FocusTabBorder           lipgloss.Border
	FocusActiveTabBorder     lipgloss.Border
	InactiveTabOnHeavyBorder lipgloss.Border

	InactiveTabStyle        lipgloss.Style
	ActiveTabStyle          lipgloss.Style
	FocusTabStyle           lipgloss.Style
	FocusActiveTabStyle     lipgloss.Style
	InactiveTabOnHeavyStyle lipgloss.Style
	TabWindowStyle          lipgloss.Style
	TabWindowHeavyStyle     lipgloss.Style
)

func init() {
	InactiveTabBorder = tabBorderWithBottom(lipgloss.RoundedBorder(), "┴", "─", "┴")
	ActiveTabBorder = tabBorderWithBottom(lipgloss.RoundedBorder(), "┘", " ", "└")

	// Heavy-border variants for focused tabs
	FocusTabBorder = tabBorderWithBottom(heavyBorder(), "┸", "─", "┸")
	FocusActiveTabBorder = tabBorderWithBottom(heavyBorder(), "┛", " ", "┗")

	// Inactive tab sitting on a heavy-bordered content window
	InactiveTabOnHeavyBorder = tabBorderWithBottom(lipgloss.RoundedBorder(), "┷", "━", "┷")
	Refresh()
}

// HeavyBorder returns a border using heavy box-drawing characters (┏━┓┃┗┛).
func HeavyBorder() lipgloss.Border {
	return heavyBorder()
}

func heavyWindowBorder() lipgloss.Border {
	b := heavyBorder()
	b.Top = ""
	b.TopLeft = ""
	b.TopRight = ""
	return b
}

func heavyBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
}

func tabBorderWithBottom(base lipgloss.Border, left, middle, right string) lipgloss.Border {
	base.BottomLeft = left
	base.Bottom = middle
	base.BottomRight = right
	return base
}
