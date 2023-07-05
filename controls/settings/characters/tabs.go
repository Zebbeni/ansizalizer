package characters

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
	windowStyle       = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m Model) drawCharTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []State{Ascii, Unicode, Custom}

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	for i, t := range tabs {
		var tabStyle lipgloss.Style
		isFirst, isLast, isActive, isMode := i == 0, i == len(tabs)-1, m.focus == t, m.mode == t

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

		if m.mode == t {
			tabStyle = activeTabStyle.Copy()
		} else {
			tabStyle = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if isFirst && isMode {
			border.BottomLeft = "│"
		} else if isFirst && !isMode {
			border.BottomLeft = "├"
		} else if isLast && isMode {
			border.BottomRight = "└"
		} else if isLast && !isMode {
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

	charButtons := m.drawCharButtons()
	doc.WriteString(windowStyle.Copy().BorderForeground(borderColor).Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(charButtons))
	return docStyle.Render(doc.String())
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
