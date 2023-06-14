package controls

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/io"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {Open: Settings, Settings: Export},
	Left:  {Export: Settings, Settings: Open},
}

func (m Model) handleOpenUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
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

		case key.Matches(msg, io.KeyMap.Enter):
			m.active = m.focus

		case key.Matches(msg, io.KeyMap.Nav):
			switch {
			case key.Matches(msg, io.KeyMap.Right):
				if next, hasNext := navMap[Right][m.focus]; hasNext {
					m.focus = next
				}
			case key.Matches(msg, io.KeyMap.Left):
				if next, hasNext := navMap[Left][m.focus]; hasNext {
					m.focus = next
				}
			}
		}
	}
	return m, nil
}
