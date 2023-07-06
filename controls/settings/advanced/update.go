package advanced

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {
		Sampling: Dithering,
	},
	Left: {
		Dithering: Sampling,
	},
	Down: {
		Sampling:  SamplingControls,
		Dithering: DitheringControls,
	},
	Up: {
		SamplingControls:  Sampling,
		DitheringControls: Dithering,
	},
}

func (m Model) handleSamplingUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.sampling, cmd = m.sampling.Update(msg)

	if m.sampling.ShouldClose {
		m.active = Menu
		m.focus = Sampling
		m.sampling.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleDitheringUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.dithering, cmd = m.dithering.Update(msg)

	if m.dithering.ShouldClose {
		m.active = Menu
		m.focus = Dithering
		m.dithering.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	}
	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	switch m.focus {
	case Sampling:
		m.activeTab = Sampling
	case Dithering:
		m.activeTab = Dithering
	case SamplingControls:
		m.active = SamplingControls
	case DitheringControls:
		m.active = DitheringControls
	}
	return m, nil
}
