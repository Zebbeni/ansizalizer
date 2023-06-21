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
	Down
	Up
)

var navMap = map[Direction]map[State]State{
	Right: {NoPalette: Load, Load: Adapt, Adapt: Lospec},
	Left:  {Lospec: Adapt, Adapt: Load, Load: NoPalette},
	Down:  {Adapt: AdaptiveControls, Load: LoadControls, Lospec: LospecControls},
	Up:    {AdaptiveControls: Adapt, LoadControls: Load, LospecControls: Lospec},
}

func (m Model) handleMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		}
	}
	return m, nil
}

func (m Model) handleAdaptiveUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Adapter, cmd = m.Adapter.Update(msg)
	if m.Adapter.ShouldUnfocus {
		m.Adapter.IsActive = true
		m.Adapter.ShouldUnfocus = false
		m.focus = Adapt
	} else if m.Adapter.ShouldClose {
		m.Adapter.IsActive = true
		m.Adapter.ShouldClose = false
		m.ShouldClose = true
	}
	return m, cmd
}

func (m Model) handlePaletteUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Loader, cmd = m.Loader.Update(msg)
	if m.Loader.ShouldUnfocus {
		m.Loader.ShouldUnfocus = false
		m.focus = Load
	}
	return m, cmd
}

func (m Model) handleLospecUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Lospec, cmd = m.Lospec.Update(msg)
	if m.Lospec.ShouldUnfocus {
		m.Lospec.IsActive = true
		m.Lospec.ShouldUnfocus = false
		m.focus = Lospec
	} else if m.Lospec.ShouldClose {
		m.Lospec.IsActive = true
		m.Lospec.ShouldClose = false
		m.ShouldClose = true
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
	if m.IsAdaptive() && len(m.Adapter.GetCurrent().Colors()) == 0 {
		return m, event.StartAdaptingCmd
	}
	return m, event.StartRenderCmd
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
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		}
	}

	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	switch m.focus {
	case Adapt:
		m.controls = Adapt
		m.selected = Adapt
	case Load:
		m.controls = Load
		m.selected = Load
	case Lospec:
		m.controls = Lospec
		m.selected = Lospec
	case NoPalette:
		m.controls = NoPalette
		m.selected = NoPalette
	case AdaptiveControls:
		m.Adapter.IsActive = true
		m.selected = Adapt
	case LoadControls:
		m.selected = Load
	case LospecControls:
		m.Lospec.IsActive = true
		m.selected = Lospec
	}
	return m, nil
}
