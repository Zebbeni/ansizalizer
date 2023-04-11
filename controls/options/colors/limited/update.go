package limited

import tea "github.com/charmbracelet/bubbletea"

func (m Model) handleBasicUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.basic, cmd = m.basic.Update(msg)
	if m.basic.ShouldFocus {
		m.basic.ShouldFocus = false
		m.ShouldUnfocus = true
	}
	return m, cmd
}
