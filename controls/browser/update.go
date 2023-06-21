package browser

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/menu"
	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleEnter() (Model, tea.Cmd) {
	return m.updateSelected()
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.currentList().Index() == 0 && key.Matches(msg, event.KeyMap.Up) {
		m.ShouldClose = true
		return m, nil
	}

	cmds := make([]tea.Cmd, 2)
	m.lists[m.listIndex()], cmds[0] = m.currentList().Update(msg)
	m, cmds[1] = m.updateActive()
	return m, tea.Batch(cmds...)
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	// remove last list if possible (go back to previous)
	if len(m.lists) > 1 {
		m.lists = m.lists[:m.listIndex()]
		return m, nil
	}

	m.ShouldClose = true
	return m, nil
}

func (m Model) updateActive() (Model, tea.Cmd) {
	itm, ok := m.currentList().SelectedItem().(item)
	if !ok {
		panic("Unexpected list item type")
	}

	if itm.isDir == false && m.ActiveFile != itm.path {
		m.ActiveFile = itm.path
		return m, event.StartRenderCmd
	}
	return m, nil
}

func (m Model) updateSelected() (Model, tea.Cmd) {
	itm, ok := m.currentList().SelectedItem().(item)
	if !ok {
		panic("Unexpected list item type")
	}

	if itm.isDir {
		m.SelectedDir = itm.path
		m = m.addListForDirectory(itm.path)
	} else {
		m.SelectedFile = itm.path
		m.ShouldClose = true
	}

	return m, nil
}

func (m Model) addListForDirectory(dir string) Model {
	newList := menu.New(getItems(m.fileExtensions, dir), m.width)

	newList.SetShowTitle(false)
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(false)
	newList.SetShowFilter(false)

	m.lists = append(m.lists, newList)
	m.SelectedDir = dir

	return m
}
