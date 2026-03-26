package characters

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = style.BgStyle().Padding(0)
	inactiveTabStyle  = style.BgStyle().Border(inactiveTabBorder, true)
	activeTabStyle    = style.BgStyle().Border(activeTabBorder, true)
	focusTabStyle     = style.ActiveTabStyle.Copy().BorderForeground(style.SelectedColor1)
	windowStyle       = style.BgStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0)
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

		isFirst := i == 0
		isLast := i == len(tabs)-1
		isActive := m.focus == t
		showControls := m.charControls == t

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

		if showControls {
			tabStyle = style.ActiveTabStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
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

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(stateNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW, extH := max(m.width-lipgloss.Width(tabBlock)-2, 0), 1

	border := lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}

	extendedStyle := style.TabWindowStyle.Copy().Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
	extended := extendedStyle.Copy().Width(extW).Height(extH).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	charButtons := m.drawCharControls()
	doc.WriteString(style.TabWindowStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Width(lipgloss.Width(row) - style.TabWindowStyle.GetHorizontalFrameSize()).Render(charButtons))
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
