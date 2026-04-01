package style

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Tab describes a single tab for rendering.
type Tab struct {
	Label    string
	Focused  bool // this tab has keyboard focus
	Active   bool // this tab's content is currently displayed
}

// RenderTabs draws a row of tabs above a content window. It handles
// single/heavy border switching for focused/active tabs and all the
// corner-character logic.
func RenderTabs(tabs []Tab, isActive bool, content string, width int) string {
	borderColor := DimmedColor2
	if isActive {
		borderColor = NormalColor1
	}

	// Determine if the content window should use heavy borders
	var anyFocused, focusedIsActive, focusIsContent bool
	for _, t := range tabs {
		if t.Focused {
			anyFocused = true
			if t.Active {
				focusedIsActive = true
			}
		}
	}
	focusIsContent = isActive && !anyFocused
	useHeavyContent := isActive && (focusedIsActive || focusIsContent)

	// Render each tab
	renderedTabs := make([]string, 0, len(tabs)+1)
	for i, t := range tabs {
		isFirst := i == 0
		isLast := i == len(tabs)-1

		fgColor := DimmedColor2
		if isActive {
			if t.Focused {
				fgColor = SelectedColor1
			} else {
				fgColor = DimmedColor1
			}
		} else if t.Focused {
			fgColor = NormalColor2
		}

		var tabStyle lipgloss.Style
		switch {
		case isActive && t.Focused && t.Active:
			tabStyle = FocusActiveTabStyle
		case isActive && t.Focused:
			tabStyle = FocusTabStyle
		case t.Active && useHeavyContent:
			tabStyle = FocusActiveTabStyle
		case t.Active:
			tabStyle = ActiveTabStyle
		case useHeavyContent:
			tabStyle = InactiveTabOnHeavyStyle
		default:
			tabStyle = InactiveTabStyle
		}

		border, _, _, _, _ := tabStyle.GetBorder()

		if useHeavyContent {
			if isFirst && t.Active {
				border.BottomLeft = "┃"
			} else if isFirst {
				border.BottomLeft = "┢"
			}
			if isLast && t.Active {
				border.BottomRight = "┗"
			} else if isLast && t.Focused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┷"
			}
		} else {
			if isFirst && t.Active {
				border.BottomLeft = "│"
			} else if isFirst && isActive && t.Focused {
				border.BottomLeft = "┞"
			} else if isFirst {
				border.BottomLeft = "├"
			}
			if isLast && t.Active {
				border.BottomRight = "└"
			} else if isLast && isActive && t.Focused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┴"
			}
		}

		tabStyle = tabStyle.Border(border).
			BorderForeground(borderColor).
			BorderBackground(ActiveTheme.Bg).
			Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(t.Label))
	}

	// Extension to fill remaining width
	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW := width - lipgloss.Width(tabBlock)
	if extW > 1 {
		var extBorder lipgloss.Border
		if useHeavyContent {
			extBorder = lipgloss.Border{BottomLeft: "━", Bottom: "━", BottomRight: "┓"}
		} else {
			extBorder = lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
		}
		extStyle := TabWindowStyle.Border(extBorder).
			BorderForeground(borderColor).
			BorderBackground(ActiveTheme.Bg).
			Padding(0)
		renderedTabs = append(renderedTabs, extStyle.Width(extW).Height(1).Render(""))
	} else {
		// Just the closing corner when there's no room for fill
		corner := "┐"
		if useHeavyContent {
			corner = "┓"
		}
		renderedTabs = append(renderedTabs, BgStyle().Foreground(borderColor).Render(corner))
	}

	// Join tabs + content (Bottom alignment so corner chars align with tab bottom borders)
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)

	winStyle := TabWindowStyle
	if useHeavyContent {
		winStyle = TabWindowHeavyStyle
	}
	contentRendered := winStyle.
		BorderForeground(borderColor).
		BorderBackground(ActiveTheme.Bg).
		Width(lipgloss.Width(row)).
		Render(content)

	var doc strings.Builder
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(contentRendered)
	return BgStyle().Padding(0).Render(doc.String())
}
