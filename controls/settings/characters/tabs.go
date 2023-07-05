package characters

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m Model) drawCharTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []State{Ascii, Unicode}

	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, m.focus == t
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(stateNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW, extH := max(m.width-lipgloss.Width(tabBlock)-4, 0), 1

	border := lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
	extendedStyle := windowStyle.Copy().Border(border).Padding(0, 0, 0, 0)
	extended := extendedStyle.Width(extW).Height(extH).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	charButtons := m.drawCharButtons()
	doc.WriteString(windowStyle.Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(charButtons))
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

func (m Model) drawModeButtons() string {
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		styleColor := normalColor
		if m.IsActive && state == m.focus {
			styleColor = focusColor
		} else if state == m.active {
			styleColor = activeColor
		}
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styleColor).
			Foreground(styleColor)
		buttons[i] = style.Copy().Width(12).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}
