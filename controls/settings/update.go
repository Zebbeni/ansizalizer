package settings

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Down Direction = iota
	Up
)

var navMap = map[Direction]map[State]State{
	Down: {Palette: Characters, Characters: Size, Size: Sampling},
	Up:   {Sampling: Size, Size: Characters, Characters: Palette},
}

func (m Model) handleSettingsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}
	return m, nil
}

func (m Model) handleColorsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Colors, cmd = m.Colors.Update(msg)

	if m.Colors.ShouldClose {
		m.active = None
		m.Colors.IsActive = false
		m.Colors.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleCharactersUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Characters, cmd = m.Characters.Update(msg)

	if m.Characters.ShouldClose {
		m.active = None
		m.Characters.IsActive = false
		m.Characters.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleSizeUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Size, cmd = m.Size.Update(msg)
	if m.Size.ShouldClose {
		m.active = None
		m.Size.IsActive = false
		m.Size.ShouldClose = false
	}
	if m.Size.ShouldUnfocus {
		return m.handleSettingsUpdate(msg)
	}
	return m, cmd
}

func (m Model) handleSamplingUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Sampling, cmd = m.Sampling.Update(msg)

	if m.Sampling.ShouldClose {
		m.active = None
		m.Sampling.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			m.active = m.focus
			switch m.active {
			case Palette:
				m.Colors.IsActive = true
			case Characters:
				m.Characters.IsActive = true
			case Size:
				m.Size.IsActive = true
			}
		case key.Matches(msg, event.KeyMap.Esc):
			m.ShouldClose = true
		case key.Matches(msg, event.KeyMap.Nav):
			switch {
			case key.Matches(msg, event.KeyMap.Up):
				if next, hasNext := navMap[Up][m.focus]; hasNext {
					m.focus = next
				}
			case key.Matches(msg, event.KeyMap.Down):
				if next, hasNext := navMap[Down][m.focus]; hasNext {
					m.focus = next
				}
			}
		}
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Esc):
		m.ShouldClose = true
	case key.Matches(msg, event.KeyMap.Nav):
		m.ShouldUnfocus = true
	}
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}
