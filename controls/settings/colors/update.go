package colors

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
		UseTrueColor: UsePalette,
	},
	Left: {
		UsePalette: UseTrueColor,
	},
	Up: {
		Palette: UsePalette,
	},
	Down: {
		UseTrueColor: Palette,
		UsePalette:   Palette,
	},
}

func (m Model) handlePaletteUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.PaletteControls, cmd = m.PaletteControls.Update(msg)

	if m.PaletteControls.ShouldClose {
		m.PaletteControls.ShouldClose = false
		m.focus = UsePalette
	}
	return m, cmd
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case UsePalette:
		m.mode = UsePalette
	case UseTrueColor:
		m.mode = UseTrueColor
	}
	return m, nil
}

func (m Model) handleEsc() (Model, tea.Cmd) {
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
	if m.mode == UseTrueColor && focus == Palette {
		return m, nil
	}

	m.focus = focus
	switch m.focus {
	case Palette:
		m.PaletteControls.IsActive = true
	}

	return m, nil
}
