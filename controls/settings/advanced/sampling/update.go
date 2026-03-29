package sampling

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	m.list.SetDelegate(NewDelegate(false))
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	selectedItem := m.list.SelectedItem().(item)
	m.Function = selectedItem.Function
	m.ShouldClose = true
	m.list.SetDelegate(NewDelegate(false))
	return m, event.StartRenderToViewCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	if key.Matches(msg, event.KeyMap.Up) && m.list.Index() == 0 {
		m.list.SetDelegate(NewDelegate(false))
		m.ShouldClose = true
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.list.SetDelegate(NewDelegate(true))
	return m, cmd
}
