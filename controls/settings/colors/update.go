package colors

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

var navMap = map[Direction]map[State]State{
	Right: {
		UseTrueColor: AdaptOn,
		UsePalette:   AdaptOn,
	},
	Left: {
		AdaptOn:  UsePalette,
		AdaptOff: UsePalette,
	},
	Up: {
		Palette: UsePalette,
	},
	Down: {
		UseTrueColor: Palette,
		UsePalette:   Palette,
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
	case UsePalette, UseTrueColor:
		if m.mode == UseTrueColor {
			m.mode = UsePalette
		} else {
			m.mode = UseTrueColor
		}
	case AdaptOn, AdaptOff:
		m.adaptToPalette = !m.adaptToPalette
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
	isPaletteMode := m.mode == UsePalette

	switch {
	case key.Matches(msg, event.KeyMap.Right):
		// Right from palette toggle → Adapt (only if palette mode shows Adapt)
		if next, hasNext := navMap[Right][m.focus]; hasNext && isPaletteMode {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if m.focus == AdaptOn || m.focus == AdaptOff {
			return m.setFocus(UsePalette)
		}
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		// From toggles row, go to palette tabs (only in palette mode)
		if m.focus == UseTrueColor || m.focus == UsePalette || m.focus == AdaptOn || m.focus == AdaptOff {
			if isPaletteMode {
				return m.setFocus(Palette)
			}
			m.IsActive = false
			m.ShouldClose = true
			return m, cmd
		}
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Tab):
		// Tab: try right first, then down
		if next, hasNext := navMap[Right][m.focus]; hasNext && isPaletteMode {
			return m.setFocus(next)
		}
		if isPaletteMode && (m.focus == UseTrueColor || m.focus == UsePalette || m.focus == AdaptOn || m.focus == AdaptOff) {
			return m.setFocus(Palette)
		}
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
	if m.mode == UseTrueColor && (focus == AdaptOn || focus == AdaptOff || focus == Palette) {
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
