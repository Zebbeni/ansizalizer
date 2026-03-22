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
	Down: {Colors: Characters, Characters: Size, Size: Adjust, Adjust: Advanced, Advanced: Animation},
	Up:   {Animation: Advanced, Advanced: Adjust, Adjust: Size, Size: Characters, Characters: Colors},
}

var navMapAlpha = map[Direction]map[State]State{
	Down: {Colors: Characters, Characters: Size, Size: Adjust, Adjust: Advanced, Advanced: Animation, Animation: Alpha},
	Up:   {Alpha: Animation, Animation: Advanced, Advanced: Adjust, Adjust: Size, Size: Characters, Characters: Colors},
}

func (m Model) handleSettingsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(keyMsg)
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

func (m Model) handleAdjustUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Adjust, cmd = m.Adjust.Update(msg)
	if m.Adjust.ShouldClose {
		m.active = None
		m.Adjust.IsActive = false
		m.Adjust.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleAdvancedUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Advanced, cmd = m.Advanced.Update(msg)

	if m.Advanced.ShouldClose {
		m.active = None
		m.Advanced.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleAnimationUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Animation, cmd = m.Animation.Update(msg)
	if m.Animation.ShouldClose {
		m.active = None
		m.Animation.IsActive = false
		m.Animation.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleAlphaUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Alpha, cmd = m.Alpha.Update(msg)

	if m.Alpha.ShouldUnfocus {
		m.active = None
		return m.handleSettingsUpdate(msg)
	}
	return m, cmd
}


func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	m.expanded[m.focus] = true
	switch m.active {
	case Colors:
		m.Colors.IsActive = true
	case Characters:
		m.Characters.IsActive = true
	case Size:
		m.Size.IsActive = true
	case Adjust:
		m.Adjust.IsActive = true
	case Advanced:
		m.Advanced.IsActive = true
	case Animation:
		m.Animation.IsActive = true
	case Alpha:
		m.Alpha.IsActive = true
	}
	return m, nil
}

func (m Model) handleExpand() (Model, tea.Cmd) {
	m.expanded[m.focus] = true
	return m, nil
}

func (m Model) handleCollapse() (Model, tea.Cmd) {
	delete(m.expanded, m.focus)
	return m, nil
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	prev := m.focus
	switch {
	case key.Matches(msg, event.KeyMap.Down):
		if m.Alpha.AlphaImage {
			if next, hasNext := navMapAlpha[Down][m.focus]; hasNext {
				m.focus = next
			}
		} else {
			if next, hasNext := navMap[Down][m.focus]; hasNext {
				m.focus = next
			}
		}
	case key.Matches(msg, event.KeyMap.Up):
		if m.Alpha.AlphaImage {
			if next, hasNext := navMapAlpha[Up][m.focus]; hasNext {
				m.focus = next
			} else {
				m.ShouldClose = true
			}
		} else {
			if next, hasNext := navMap[Up][m.focus]; hasNext {
				m.focus = next
			} else {
				m.ShouldClose = true
			}
		}
	}
	if m.focus != prev {
		m.resetSubMenu(prev)
	}
	return m, nil
}

func (m *Model) resetSubMenu(state State) {
	switch state {
	case Colors:
		m.Colors.ResetFocus()
	case Characters:
		m.Characters.ResetFocus()
	case Size:
		m.Size.ResetFocus()
	case Adjust:
		m.Adjust.ResetFocus()
	case Advanced:
		m.Advanced.ResetFocus()
	case Animation:
		m.Animation.ResetFocus()
	case Alpha:
		m.Alpha.ResetFocus()
	}
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleEnter()
	case key.Matches(msg, event.KeyMap.Expand):
		return m.handleExpand()
	case key.Matches(msg, event.KeyMap.Collapse):
		return m.handleCollapse()
	case key.Matches(msg, event.KeyMap.Nav):
		return m.handleNav(msg)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.handleEsc()
	}
	return m, cmd
}
