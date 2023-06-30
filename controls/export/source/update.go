package source

import (
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

var (
	navMap = map[Direction]map[State]State{
		Right: {ExpFile: ExpDirectory},
		Left:  {ExpDirectory: ExpFile},
		Down:  {ExpFile: SrcInput, ExpDirectory: SrcInput},
		Up:    {SrcInput: ExpFile},
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
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

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case ExpFile:
		m.doExportDirectory = false
	case ExpDirectory:
		m.doExportDirectory = true
	case SubDirYes:
		m.doExportSubDirectories = true
	case SubDirsNo:
		m.doExportSubDirectories = false
	case SrcInput:
		m.focus = SrcBrowser
	}
	return m, nil
}

func (m Model) handleSrcBrowserUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.SourceBrowser, cmd = m.SourceBrowser.Update(msg)
	m.sourceFilepath = m.SourceBrowser.ActiveFile

	if m.SourceBrowser.ShouldClose {
		m.focus = SrcInput
		m.SourceBrowser.ShouldClose = false
	}
	return m, cmd
}
