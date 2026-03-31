package characters

import (
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawCharTabs() string {
	tabs := []style.Tab{
		{Label: stateNames[Unicode], Focused: m.focus == Unicode, Active: m.charControls == Unicode},
		{Label: stateNames[Ascii], Focused: m.focus == Ascii, Active: m.charControls == Ascii},
		{Label: stateNames[Custom], Focused: m.focus == Custom, Active: m.charControls == Custom},
	}
	tabActive := m.IsActive
	return style.RenderTabs(tabs, tabActive, m.drawCharControls(), m.width)
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
