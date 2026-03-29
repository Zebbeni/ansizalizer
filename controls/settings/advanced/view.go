package advanced

import (
	"github.com/Zebbeni/ansizalizer/style"
)

var stateNames = map[State]string{Sampling: "Sampling", Dithering: "Dithering"}

func (m Model) drawTabs() string {
	var tabs []style.Tab
	if m.ShowDithering {
		tabs = append(tabs, style.Tab{Label: stateNames[Dithering], Focused: m.focus == Dithering, Active: m.activeTab == Dithering})
	}
	tabs = append(tabs, style.Tab{Label: stateNames[Sampling], Focused: m.focus == Sampling, Active: m.activeTab == Sampling})

	return style.RenderTabs(tabs, m.IsActive, m.drawTabContent(), m.width)
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
