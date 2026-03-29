package characters

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawCharTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []State{Ascii, Unicode, Custom}

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	// Determine if content window should use double border:
	// focused tab == active tab, or focus is inside tab content
	focusIsTab := m.focus == Ascii || m.focus == Unicode || m.focus == Custom
	focusedTabIsActive := focusIsTab && m.focus == m.charControls
	focusIsContent := m.IsActive && !focusIsTab
	useDoubleContent := m.IsActive && (focusedTabIsActive || focusIsContent)

	for i, t := range tabs {
		isFirst := i == 0
		isLast := i == len(tabs)-1
		isFocused := m.focus == t
		showControls := m.charControls == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isFocused {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if isFocused {
				fgColor = style.NormalColor2
			}
		}

		// Pick tab style based on focus and active state
		var tabStyle lipgloss.Style
		if m.IsActive && isFocused && showControls {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if m.IsActive && isFocused {
			tabStyle = style.FocusTabStyle.Copy()
		} else if showControls && useDoubleContent {
			// Active tab gets double to connect cleanly with double content
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if showControls {
			tabStyle = style.ActiveTabStyle.Copy()
		} else if useDoubleContent {
			tabStyle = style.InactiveTabOnHeavyStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()

		// Override corner characters based on position
		if useDoubleContent {
			if isFirst && showControls {
				border.BottomLeft = "┃"
			} else if isFirst {
				border.BottomLeft = "┢"
			}
			if isLast && showControls {
				border.BottomRight = "┗"
			} else if isLast && isFocused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┷"
			}
		} else {
			if isFirst && showControls {
				border.BottomLeft = "│"
			} else if isFirst && m.IsActive && isFocused {
				border.BottomLeft = "┞"
			} else if isFirst {
				border.BottomLeft = "├"
			}
			if isLast && showControls {
				border.BottomRight = "└"
			} else if isLast && m.IsActive && isFocused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┴"
			}
		}

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(stateNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW, extH := max(m.width-lipgloss.Width(tabBlock)-2, 0), 1

	var extBorder lipgloss.Border
	if useDoubleContent {
		extBorder = lipgloss.Border{BottomLeft: "━", Bottom: "━", BottomRight: "┓"}
	} else {
		extBorder = lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
	}

	extendedStyle := style.TabWindowStyle.Border(extBorder).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
	extended := extendedStyle.Width(extW).Height(extH).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	winStyle := style.TabWindowStyle
	if useDoubleContent {
		winStyle = style.TabWindowHeavyStyle
	}

	charButtons := m.drawCharControls()
	doc.WriteString(winStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).AlignHorizontal(lipgloss.Center).Width(lipgloss.Width(row) - winStyle.GetHorizontalFrameSize()).Render(charButtons))
	return style.BgStyle().Padding(0).Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
