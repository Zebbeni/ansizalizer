package controls

import (
	"os"

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
	Right: {Browse: Settings, Settings: Export},
	Left:  {Export: Settings, Settings: Browse},
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
			m.active = m.focus
			m.showing = m.focus
			if m.focus == Browse {
				m.FileBrowser.SetActive(true)
			}

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
				// Re-enter the tab whose content is showing;
				// focused-but-inactive tabs do not navigate on down.
				if m.focus == m.showing {
					m.active = m.focus
					if m.focus == Browse {
						m.FileBrowser.SetActive(true)
					}
				}
			}

		case key.Matches(msg, event.KeyMap.Esc):
			// Quit program if top-level menu is active and escape pressed
			tea.Quit()
			os.Exit(0)
		}
	}
	return m, nil
}
