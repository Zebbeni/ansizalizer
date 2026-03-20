package alpha

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

var (
	navMap = map[Direction]map[State]State{
		Right: {AlphaYes: AlphaNo, TrimAlphaYes: TrimAlphaNo},
		Left:  {AlphaNo: AlphaYes, TrimAlphaNo: TrimAlphaYes},
		Up:    {TrimAlphaYes: AlphaYes, TrimAlphaNo: AlphaNo},
		Down:  {AlphaYes: TrimAlphaYes, AlphaNo: TrimAlphaNo},
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldUnfocus = true
	m.IsActive = false
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldUnfocus = true
		}
	}
	return m, nil
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.IsActive = true
	m.ShouldUnfocus = false
	m.focus = focus
	if m.focus == AlphaYes {
		m.useAlpha = true
	}
	if m.focus == AlphaNo {
		m.useAlpha = false
	}
	if m.focus == TrimAlphaYes {
		m.trimAlpha = true
	}
	if m.focus == TrimAlphaNo {
		m.trimAlpha = false
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case AlphaYes:
		m.useAlpha = true
	case AlphaNo:
		m.useAlpha = false
		m.trimAlpha = false
	case TrimAlphaYes:
		m.trimAlpha = true
	case TrimAlphaNo:
		m.trimAlpha = false
	}
	return m, event.StartRenderToViewCmd
}
