package export

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
	Down: {Source: Destination, Destination: Process},
	Up:   {Destination: Source, Process: Destination},
}

func (m Model) handleSourceUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Source, cmd = m.Source.Update(msg)

	if m.Source.ShouldClose {
		m.active = None
		m.Source.ShouldClose = false
	}
	if m.Source.ShouldUnfocus {
		return m.handleMenuUpdate(msg)
	}
	return m, cmd
}

func (m Model) handleDestinationUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Destination, cmd = m.Destination.Update(msg)

	if m.Destination.ShouldClose {
		m.active = None
		m.Destination.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	switch m.active {
	case Source:
		m.Source.IsActive = true
	case Destination:
		m.Destination.IsActive = true
	case Process:
		return m.handleProcess()
	}
	return m, nil
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			m.focus = next
		} else {
			m.ShouldClose = true
		}
	}
	return m, nil
}

func (m Model) handleMenuUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(keyMsg)
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleEnter()
	case key.Matches(msg, event.KeyMap.Nav):
		return m.handleNav(msg)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.handleEsc()
	}
	return m, cmd
}

func (m Model) handleProcess() (Model, tea.Cmd) {
	sourcePath, isDir, useSubDirs := m.Source.GetSelected()
	destinationPath := m.Destination.GetSelected()
	return m, event.BuildStartExportCmd(event.StartExportMsg{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		IsDir:           isDir,
		UseSubDirs:      useSubDirs,
	})
}

func (m Model) GetDestination() (path string) {
	return m.Destination.GetSelected()
}
