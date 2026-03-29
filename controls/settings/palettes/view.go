package palettes

import (
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
	tabs := make([]style.Tab, len(stateOrder))
	for i, t := range stateOrder {
		tabs[i] = style.Tab{Label: stateNames[t], Focused: m.focus == t, Active: m.controls == t}
	}
	return style.RenderTabs(tabs, m.IsActive, m.drawTabContent(), m.width)
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
