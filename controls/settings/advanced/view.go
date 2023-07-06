package advanced

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0)
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true)
	activeTabStyle    = lipgloss.NewStyle().Border(activeTabBorder, true)
	focusTabStyle     = activeTabStyle.Copy().BorderForeground(style.SelectedColor1)
	windowStyle       = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0)
	stateNames        = map[State]string{Sampling: "Sampling", Dithering: "Dithering"}
)

func (m Model) drawTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []State{Sampling, Dithering}

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	for i, t := range tabs {
		var tabStyle lipgloss.Style
		isFirst, isLast, isActive, isActiveTab := i == 0, i == len(tabs)-1, m.focus == t, m.activeTab == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isActive {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if isActive {
				fgColor = style.NormalColor2
			}
		}

		if m.activeTab == t {
			tabStyle = activeTabStyle.Copy()
		} else {
			tabStyle = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if isFirst && isActiveTab {
			border.BottomLeft = "│"
		} else if isFirst && !isActiveTab {
			border.BottomLeft = "├"
		} else if isLast && isActiveTab {
			border.BottomRight = "└"
		} else if isLast && !isActiveTab {
			border.BottomRight = "┴"
		}

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(stateNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW, extH := max(m.width-lipgloss.Width(tabBlock)-2, 0), 1

	border := lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}

	extendedStyle := windowStyle.Copy().Border(border).BorderForeground(borderColor).Padding(0)
	extended := extendedStyle.Copy().Width(extW).Height(extH).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	content := m.drawTabContent()
	doc.WriteString(windowStyle.Copy().BorderForeground(borderColor).Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(content))
	return docStyle.Render(doc.String())
}

func (m Model) drawTabContent() string {
	switch m.activeTab {
	case Sampling:
		return m.sampling.View()
	case Dithering:
		return "Dithering Content"
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
