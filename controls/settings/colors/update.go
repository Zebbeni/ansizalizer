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
	Up
)

var navMap = map[Direction]map[State]State{
	Right: {NoPalette: Load, Load: Create, Create: Lospec},
	Left:  {Lospec: Create, Create: Load, Load: NoPalette},
	Down:  {Create: CreateControls, Load: LoadControls},
	Up:    {CreateControls: Create},
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
	m.Creator, cmd = m.Creator.Update(msg)
	if m.Creator.ShouldUnfocus {
		m.Creator.IsActive = true
		m.Creator.ShouldUnfocus = false
		m.focus = Create
		return m, cmd
	} else if m.Creator.ShouldClose {
		m.Creator.IsActive = true
		m.Creator.ShouldClose = false
		m.ShouldClose = true
		return m, cmd
	}
	return m, cmd
}

func (m Model) handlePaletteUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Loader, cmd = m.Loader.Update(msg)
	if m.Loader.ShouldUnfocus {
		m.Loader.ShouldUnfocus = false
		m.focus = Load
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
	// Kick off a new palette generation before rendering if not done yet.
	// Allow the app to trigger a render when the generation is complete.
	if m.IsAdaptive() && len(m.Creator.Palette) == 0 {
		return m, io.BuildAdaptingCmd()
	}
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
	case key.Matches(msg, io.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		}
	}

	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	switch m.focus {
	case Create:
		m.controls = Create
		m.selected = Create
	case Load:
		m.controls = Load
		m.selected = Load
	case NoPalette:
		m.controls = NoPalette
		m.selected = NoPalette
	case CreateControls:
		m.Creator.IsActive = true
		m.selected = Create
	case LoadControls:
		m.selected = Load
	}
	return m, nil
}
