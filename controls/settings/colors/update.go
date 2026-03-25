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
		AdaptOn:      AdaptOff,
	},
	Left: {
		UsePalette: UseTrueColor,
		AdaptOff:   AdaptOn,
	},
	Up: {
		AdaptOn:  UsePalette,
		AdaptOff: UsePalette,
		Palette:  AdaptOn,
	},
	Down: {
		UseTrueColor: AdaptOn,
		UsePalette:   AdaptOn,
		AdaptOn:      Palette,
		AdaptOff:     Palette,
	},
}

func (m Model) handlePaletteUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.PaletteControls, cmd = m.PaletteControls.Update(msg)

	if m.PaletteControls.ShouldClose {
		m.PaletteControls.IsActive = false
		m.PaletteControls.ShouldClose = false
		if m.PaletteControls.IsLospecFocused() {
			m.focus = AdaptOff
		} else {
			m.focus = AdaptOn
		}
	}
	return m, cmd
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case UsePalette:
		m.mode = UsePalette
	case UseTrueColor:
		m.mode = UseTrueColor
	case AdaptOn:
		m.adaptToPalette = true
	case AdaptOff:
		m.adaptToPalette = false
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.IsActive = false
	m.ShouldClose = true
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
		m.IsActive = false
		m.ShouldClose = true
		return m, nil
	}

	m.focus = focus
	if m.focus == Palette {
		m.PaletteControls.IsActive = true
	} else {
		m.PaletteControls.IsActive = false
	}

	return m, nil
}
