package settings

import (
	"path/filepath"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Down Direction = iota
	Up
)

func (m Model) visibleStates() []State {
	states := []State{Colors, Characters, Size, Adjust, TextStyle, Advanced}
	if m.Animation.AnimatedImage {
		states = append(states, Animation)
	}
	if m.Alpha.AlphaImage {
		states = append(states, Alpha)
	}
	return states
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

func (m Model) handleTextStyleUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.TextStyle, cmd = m.TextStyle.Update(msg)
	if m.TextStyle.ShouldClose {
		m.active = None
		m.TextStyle.IsActive = false
		m.TextStyle.ShouldClose = false
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

func (m Model) handleThemeUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Theme, cmd = m.Theme.Update(msg)

	if m.Theme.ShouldClose {
		m.active = None
		m.Theme.IsActive = false
		m.Theme.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleSaveLoadUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.SaveLoad, cmd = m.SaveLoad.Update(msg)

	if m.SaveLoad.ShouldClose {
		// Check if a save or load was requested before closing
		if m.SaveLoad.SavePath != "" {
			cfg := m.ExportConfig()
			err := SaveConfig(cfg, m.SaveLoad.SavePath)
			filename := filepath.Base(m.SaveLoad.SavePath)
			m.SaveLoad.SavePath = ""
			m.SaveLoad.ShouldClose = false
			if err != nil {
				m.SaveLoad.SetSaveResult("Save Failed")
				return m, event.BuildDisplayCmd("error saving settings")
			}
			m.SaveLoad.SetSaveResult("Saved to " + filename)
			m.SaveLoad.SetStatus("Saved " + filename)
			return m, event.BuildDisplayCmd("settings saved")
		}
		if m.SaveLoad.LoadPath != "" {
			cfg, err := LoadConfig(m.SaveLoad.LoadPath)
			filename := filepath.Base(m.SaveLoad.LoadPath)
			m.SaveLoad.LoadPath = ""
			m.SaveLoad.ShouldClose = false
			if err != nil {
				m.SaveLoad.SetLoadResult("Error Loading " + filename)
				return m, event.BuildDisplayCmd("error loading settings")
			}
			m.ApplyConfig(cfg)
			m.SaveLoad.SetLoadResult("Settings File: " + filename)
			m.SaveLoad.SetStatus("Loaded " + filename)
			return m, tea.Batch(event.StartRenderToViewCmd, event.BuildDisplayCmd("settings loaded"))
		}
		m.active = None
		m.SaveLoad.IsActive = false
		m.SaveLoad.ShouldClose = false
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
	m.expanded = map[State]bool{m.focus: true}
	if m.active != SaveLoad {
		m.SaveLoad.ClearStatus()
	}
	switch m.active {
	case Colors:
		m.Colors.IsActive = true
	case Characters:
		m.Characters.IsActive = true
	case Size:
		m.Size.IsActive = true
	case Adjust:
		m.Adjust.IsActive = true
	case TextStyle:
		m.TextStyle.IsActive = true
	case Advanced:
		m.Advanced.IsActive = true
	case Animation:
		m.Animation.IsActive = true
	case Alpha:
		m.Alpha.IsActive = true
	case Theme:
		m.Theme.IsActive = true
	case SaveLoad:
		m.SaveLoad.IsActive = true
	}
	return m, nil
}

func (m Model) handleExpand() (Model, tea.Cmd) {
	m.expanded = map[State]bool{m.focus: true}
	return m, nil
}

func (m Model) handleCollapse() (Model, tea.Cmd) {
	delete(m.expanded, m.focus)
	return m, nil
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	if m.expanded[m.focus] {
		delete(m.expanded, m.focus)
		return m, nil
	}
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	prev := m.focus
	states := m.visibleStates()

	switch {
	case key.Matches(msg, event.KeyMap.Down):
		for i, s := range states {
			if s == m.focus && i+1 < len(states) {
				m.focus = states[i+1]
				break
			}
		}
	case key.Matches(msg, event.KeyMap.Tab):
		for i, s := range states {
			if s == m.focus && i+1 < len(states) {
				m.focus = states[i+1]
				break
			}
		}
	case key.Matches(msg, event.KeyMap.Up):
		for i, s := range states {
			if s == m.focus {
				if i == 0 {
					m.ShouldClose = true
				} else {
					m.focus = states[i-1]
				}
				break
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
	case TextStyle:
		m.TextStyle.ResetFocus()
	case Advanced:
		m.Advanced.ResetFocus()
	case Animation:
		m.Animation.ResetFocus()
	case Alpha:
		m.Alpha.ResetFocus()
	case Theme:
		m.Theme.ResetFocus()
	case SaveLoad:
		m.SaveLoad.ResetFocus()
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
