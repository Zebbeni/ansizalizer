package basic

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/io"
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldFocus = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	selectedItem := m.menu.SelectedItem().(item)

	if selectedItem.name == m.name {
		m.ShouldFocus = true
	}
	m.name = selectedItem.name
	m.palette = selectedItem.palette

	return m, io.StartRenderCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	// if at the top of the menu, set flag to hand focus back to parent.
	if m.menu.Index() == 0 && key.Matches(msg, io.KeyMap.Up) {
		m.ShouldFocus = true
		return m, nil
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}
