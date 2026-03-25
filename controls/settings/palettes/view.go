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

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true)
	activeTabStyle    = lipgloss.NewStyle().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0)
)

func (m Model) drawTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	for i, t := range stateOrder {
		var tabStyle lipgloss.Style

		isFirst := i == 0
		isLast := i == len(stateOrder)-1
		isActive := m.focus == t
		showControls := m.controls == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isActive {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if showControls {
				fgColor = style.NormalColor2
			}
		}

		if showControls {
			tabStyle = activeTabStyle.Copy()
		} else {
			tabStyle = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if isFirst && showControls {
			border.BottomLeft = "│"
		} else if isFirst && !showControls {
			border.BottomLeft = "├"
		} else if isLast && showControls {
			border.BottomRight = "└"
		} else if isLast && !showControls {
			border.BottomRight = "┴"
		}

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(stateNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW := m.width - lipgloss.Width(tabBlock) - 4

	if extW > 0 {
		border := lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
		extendedStyle := windowStyle.Copy().Border(border).BorderForeground(borderColor).Padding(0)
		extended := extendedStyle.Copy().Width(extW).Height(1).Render("")
		renderedTabs = append(renderedTabs, extended)
	} else {
		// Just add the closing corner, aligned to the bottom of the tab row
		corner := lipgloss.NewStyle().Foreground(borderColor).Render("┐")
		renderedTabs = append(renderedTabs, corner)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	controls := m.drawTabContent()
	doc.WriteString(windowStyle.Copy().BorderForeground(borderColor).Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(controls))
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
