package dithering

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
		DitherOn:     DitherOff,
		SerpentineOn: SerpentineOff,
	},
	Left: {
		DitherOff:     DitherOn,
		SerpentineOff: SerpentineOn,
	},
	Down: {
		DitherOn:      SerpentineOn,
		DitherOff:     SerpentineOff,
		SerpentineOn:  Matrix,
		SerpentineOff: Matrix,
	},
	Up: {
		SerpentineOn:  DitherOn,
		SerpentineOff: DitherOff,
		Matrix:        SerpentineOn,
	},
}

func (m Model) handleMatrixListUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Up) && m.list.Index() == 0:
			return m.handleNav(keyMsg)
		case key.Matches(keyMsg, event.KeyMap.Esc):
		case key.Matches(keyMsg, event.KeyMap.Enter):
			var cmd tea.Cmd
			m, cmd = m.setFocus(navMap[Up][Matrix])
			return m, tea.Batch(cmd, event.StartRenderToViewCmd)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case DitherOn:
		m.doDithering = true
	case DitherOff:
		m.doDithering = false
	case SerpentineOn:
		m.doSerpentine = true
	case SerpentineOff:
		m.doSerpentine = false
	}
	return m, event.StartRenderToViewCmd
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
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.ShouldClose = true
		}
	}
	return m, cmd
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleEnter()
	case key.Matches(msg, event.KeyMap.Nav):
		return m.handleNav(msg)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.handleEsc()
	}
	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	if focus != Matrix {
		m.list.SetDelegate(NewDelegate(false))
	} else {
		m.list.SetDelegate(NewDelegate(true))
	}

	return m, nil
}
