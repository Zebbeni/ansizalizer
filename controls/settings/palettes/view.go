package palettes

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder = []State{Lospec, Load, Adapt}
	stateNames = map[State]string{
		Lospec: "Lospec",
		Load:   "From File",
		Adapt:  "Sample",
	}

)

func (m Model) drawTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	focusIsTab := m.focus == Lospec || m.focus == Load || m.focus == Adapt
	focusedTabIsActive := focusIsTab && m.focus == m.controls
	focusIsContent := m.IsActive && !focusIsTab
	useDoubleContent := m.IsActive && (focusedTabIsActive || focusIsContent)

	for i, t := range stateOrder {
		isFirst := i == 0
		isLast := i == len(stateOrder)-1
		isFocused := m.focus == t
		showControls := m.controls == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isFocused {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if showControls {
				fgColor = style.NormalColor2
			}
		}

		var tabStyle lipgloss.Style
		if m.IsActive && isFocused && showControls {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if m.IsActive && isFocused {
			tabStyle = style.FocusTabStyle.Copy()
		} else if showControls && useDoubleContent {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if showControls {
			tabStyle = style.ActiveTabStyle.Copy()
		} else if useDoubleContent {
			tabStyle = style.InactiveTabOnHeavyStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
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
	extW := m.width - lipgloss.Width(tabBlock) - 4

	if extW > 0 {
		var extBorder lipgloss.Border
		if useDoubleContent {
			extBorder = lipgloss.Border{BottomLeft: "━", Bottom: "━", BottomRight: "┓"}
		} else {
			extBorder = lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
		}
		extendedStyle := style.TabWindowStyle.Copy().Border(extBorder).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
		extended := extendedStyle.Copy().Width(extW).Height(1).Render("")
		renderedTabs = append(renderedTabs, extended)
	} else {
		corner := "┐"
		if useDoubleContent {
			corner = "┓"
		}
		renderedTabs = append(renderedTabs, style.BgStyle().Foreground(borderColor).Render(corner))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	winStyle := style.TabWindowStyle
	if useDoubleContent {
		winStyle = style.TabWindowHeavyStyle
	}
	controls := m.drawTabContent()
	doc.WriteString(winStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Width(lipgloss.Width(row) - winStyle.GetHorizontalFrameSize()).Render(controls))
	return doc.String()
}

func (m Model) drawTabContent() string {
	switch m.controls {
	case Adapt:
		return m.Adapter.View()
	case Load:
		return m.Loader.View()
	case Lospec:
		return m.Lospec.View()
	}
	return ""
}


func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
