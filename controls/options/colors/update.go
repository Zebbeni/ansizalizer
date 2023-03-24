package colors

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/io"
)

type Direction int

const (
	Left Direction = iota
	Right
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {TrueColor: Adaptive, Adaptive: Paletted},
	Left:  {Paletted: Adaptive, Adaptive: TrueColor},
	Down:  {Adaptive: AdaptiveControls, Paletted: PalettedControls},
}

func (m Model) handleMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, io.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, io.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, io.KeyMap.Nav):
			return m.handleNav(msg)
		}
	}
	return m, nil
}

func (m Model) handleAdaptiveUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Adaptive, cmd = m.Adaptive.Update(msg)
	if m.Adaptive.ShouldUnfocus {
		m.Adaptive.ShouldUnfocus = false
		m.focus = Adaptive
		return m, cmd
	}
	return m, cmd
}

func (m Model) handlePaletteUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Palette, cmd = m.Palette.Update(msg)
	if m.Palette.ShouldUnfocus {
		m.Palette.ShouldUnfocus = false
		m.focus = Paletted
		return m, cmd
	}
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.selected = m.focus
	return m, io.StartRenderCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, io.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, io.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, io.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		}
	}

	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	switch m.focus {
	case Adaptive:
		m.controls = Adaptive
	case Paletted:
		m.controls = Paletted
	case TrueColor:
		m.controls = TrueColor
	}
	return m, nil
}
