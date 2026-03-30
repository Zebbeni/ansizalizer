package characters

import (
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawCharTabs() string {
	tabs := []style.Tab{
		{Label: stateNames[Ascii], Focused: m.focus == Ascii, Active: m.charControls == Ascii},
		{Label: stateNames[Unicode], Focused: m.focus == Unicode, Active: m.charControls == Unicode},
		{Label: stateNames[Custom], Focused: m.focus == Custom, Active: m.charControls == Custom},
	}
	// Only pass isActive when focus is on a tab or inside tab content
	// (not on ColorBgOn/Off which are above the tabs)
	tabActive := m.IsActive && m.focus != ColorBgOn && m.focus != ColorBgOff
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
