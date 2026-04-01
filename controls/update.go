package controls

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/controls/filemenu"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {FileMenu: Browse, Browse: Settings},
	Left:  {Settings: Browse, Browse: FileMenu},
}

func (m Model) handleFileMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	m.FileDropdown, _ = m.FileDropdown.Update(msg)

	if result := m.FileDropdown.SelectedResult; result != nil {
		m.FileDropdown.SelectedResult = nil
		return m.handleFileMenuResult(*result)
	}

	if m.FileDropdown.ShouldNavRight {
		m.FileDropdown.ShouldNavRight = false
		m.focus = Browse
		m.active = Menu
		return m, nil
	}

	if m.FileDropdown.ShouldClose {
		m.FileDropdown.ShouldClose = false
		m.focus = FileMenu
		m.active = Menu
		return m, nil
	}

	if !m.FileDropdown.Open {
		m.active = Menu
	}
	return m, nil
}

func (m Model) handleFileMenuResult(result filemenu.Result) (Model, tea.Cmd) {
	switch result.Action {
	case filemenu.LoadSettings:
		m.OpenModal = &result
		return m, nil
	case filemenu.SaveSettings:
		m.OpenModal = &result
		return m, nil
	case filemenu.Theme:
		if result.ThemeName != "" {
			return m, event.BuildCheckThemeCmd(style.PaletteColors, result.ThemeName)
		}
	case filemenu.Export:
		m.OpenModal = &result
	case filemenu.Preferences:
		m.OpenModal = &result
	case filemenu.Quit:
		m.ShouldQuit = true
		return m, nil
	}
	return m, nil
}

func (m Model) handleOpenUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
		m.FileBrowser.SetActive(false)
		m.active = Menu
	}
	return m, cmd
}

func (m Model) handleSettingsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Settings, cmd = m.Settings.Update(msg)

	if m.Settings.ShouldClose {
		m.Settings.ShouldClose = false
		m.Settings.IsActive = false
		m.active = Menu
	}

	return m, cmd
}

func (m Model) handleExportUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Export, cmd = m.Export.Update(msg)

	if m.Export.ShouldClose {
		m.Export.ShouldClose = false
		m.active = Menu
	}

	return m, cmd
}

func (m Model) handleMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	m.active = Menu
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			if m.focus == FileMenu {
				m.FileDropdown.HasSelectedFile = m.FileBrowser.SelectedFile != ""
				m.FileDropdown.Toggle()
				return m, nil
			}
			m.active = m.focus
			m.showing = m.focus
			m.activateShowing()

		case key.Matches(msg, event.KeyMap.Nav):
			switch {
			case key.Matches(msg, event.KeyMap.Right):
				if next, hasNext := navMap[Right][m.focus]; hasNext {
					m.focus = next
				}
			case key.Matches(msg, event.KeyMap.Left):
				if next, hasNext := navMap[Left][m.focus]; hasNext {
					m.focus = next
				}
			case key.Matches(msg, event.KeyMap.Down):
				m.enterShowing()
			case key.Matches(msg, event.KeyMap.Tab):
				// Try Right first
				if next, hasNext := navMap[Right][m.focus]; hasNext {
					m.focus = next
					return m, nil
				}
				// Fall through to Down logic
				m.enterShowing()
			}

		case key.Matches(msg, event.KeyMap.Esc):
			m.ShouldQuit = true
		}
	}
	return m, nil
}

// activateShowing sets the active-state flags for whichever content is showing.
func (m *Model) activateShowing() {
	if m.showing == Browse {
		m.FileBrowser.SetActive(true)
	}
	if m.showing == Settings {
		m.Settings.IsActive = true
	}
}

// enterShowing navigates from the button bar into the showing content.
func (m *Model) enterShowing() {
	if m.showing == Browse || m.showing == Settings {
		m.active = m.showing
		m.focus = m.showing
		m.activateShowing()
	}
}
