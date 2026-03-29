package advanced

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var stateNames = map[State]string{Sampling: "Sampling", Dithering: "Dithering"}

func (m Model) drawTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	var tabs []State
	if m.ShowDithering {
		tabs = append(tabs, Dithering)
	}
	tabs = append(tabs, Sampling)

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	focusIsTab := m.focus == Sampling || m.focus == Dithering
	focusedTabIsActive := focusIsTab && m.focus == m.activeTab
	focusIsContent := m.IsActive && !focusIsTab
	useDoubleContent := m.IsActive && (focusedTabIsActive || focusIsContent)

	for i, t := range tabs {
		isFirst, isLast := i == 0, i == len(tabs)-1
		isFocused := m.focus == t
		isActiveTab := m.activeTab == t

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

		var tabStyle lipgloss.Style
		if m.IsActive && isFocused && isActiveTab {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if m.IsActive && isFocused {
			tabStyle = style.FocusTabStyle.Copy()
		} else if isActiveTab && useDoubleContent {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if isActiveTab {
			tabStyle = style.ActiveTabStyle.Copy()
		} else if useDoubleContent {
			tabStyle = style.InactiveTabOnHeavyStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if useDoubleContent {
			if isFirst && isActiveTab {
				border.BottomLeft = "┃"
			} else if isFirst {
				border.BottomLeft = "┢"
			}
			if isLast && isActiveTab {
				border.BottomRight = "┗"
			} else if isLast && isFocused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┷"
			}
		} else {
			if isFirst && isActiveTab {
				border.BottomLeft = "│"
			} else if isFirst && m.IsActive && isFocused {
				border.BottomLeft = "┞"
			} else if isFirst {
				border.BottomLeft = "├"
			}
			if isLast && isActiveTab {
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

	extendedStyle := style.TabWindowStyle.Copy().Border(extBorder).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
	extended := extendedStyle.Copy().Width(extW).Height(extH).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	winStyle := style.TabWindowStyle
	if useDoubleContent {
		winStyle = style.TabWindowHeavyStyle
	}

	content := m.drawTabContent()
	doc.WriteString(winStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Width(lipgloss.Width(row) - winStyle.GetHorizontalFrameSize()).Render(content))
	return style.BgStyle().Padding(0).Render(doc.String())
}

func (m Model) drawTabContent() string {
	switch m.activeTab {
	case Sampling:
		return m.sampling.View()
	case Dithering:
		return m.dithering.View()
	}
	return ""
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
