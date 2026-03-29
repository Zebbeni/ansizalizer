package animation

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	if m.active == m.focus {
		switch m.active {
		case DelayForm:
			m.delayInput.Blur()
			m.active = None
		}
		return m, event.StartRenderToViewCmd
	}

	m.active = m.focus
	switch m.active {
	case DelayForm:
		m.delayInput.Focus()
		m.delayInput.CursorEnd()
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleDelayUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.delayInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.delayInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.delayInput, cmd = m.delayInput.Update(msg)
	return m, cmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.active {
	case DelayForm:
		m.delayInput.Blur()
		m.active = None
	}

	switch {
	case key.Matches(msg, event.KeyMap.Up):
		m.ShouldClose = true
	}
	return m, nil
}
