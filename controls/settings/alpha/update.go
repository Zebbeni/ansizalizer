package alpha

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var (
	navMap = map[Direction]map[State]State{
		Down: {ThresholdForm: UseAlpha, UseAlpha: TrimAlpha},
		Up:   {UseAlpha: ThresholdForm, TrimAlpha: UseAlpha},
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	m.IsActive = false
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Down), key.Matches(msg, event.KeyMap.Tab):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			// Skip TrimAlpha if alpha is disabled
			if next == TrimAlpha && !m.useAlpha {
				m.IsActive = false
				m.ShouldClose = true
				return m, nil
			}
			return m.setFocus(next)
		}
		// Past last visible item
		m.IsActive = false
		m.ShouldClose = true
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		}
		m.IsActive = false
		m.ShouldClose = true
	}
	return m, nil
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.IsActive = true
	m.ShouldClose = false
	m.focus = focus
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case UseAlpha:
		m.useAlpha = !m.useAlpha
		if !m.useAlpha {
			m.trimAlpha = false
		}
	case TrimAlpha:
		m.trimAlpha = !m.trimAlpha
	case ThresholdForm:
		m.thresholdInput.Focus()
		return m, nil
	}
	return m, event.StartRenderToViewCmd
}
