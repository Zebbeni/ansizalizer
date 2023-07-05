package destination

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Up Direction = iota
	Down
)

var (
	navMap = map[Direction]map[State]State{
		Down: {Input: Browser},
		Up:   {Browser: Input},
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	m.IsActive = false
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

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case Input:
		m.focus = Browser
	}
	return m, nil
}

func (m Model) handleDstBrowserUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Browser, cmd = m.Browser.Update(msg)
	m.selectedDir = m.Browser.SelectedDir

	if m.Browser.ShouldClose {
		m.focus = Input
		m.Browser.ShouldClose = false
	}
	return m, cmd
}
