package browser

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/menu"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
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

	newList.Styles.Title = style.DimmedTitle
	newList.Styles.TitleBar = newList.Styles.TitleBar.Copy().Padding(1, 0, 1, 0).Width(m.width).AlignHorizontal(lipgloss.Center)
	newList.SetShowTitle(true)
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(false)
	newList.SetShowFilter(false)
	newList.SetWidth(m.width)
	title := fmt.Sprintf("Select %s files", m.fileExtensionsString())
	newList.Title = title

	m.lists = append(m.lists, newList)
	m.SelectedDir = dir

	return m
}

func (m Model) fileExtensionsString() string {
	toReturn := ""
	i := 0
	for ; i < len(m.fileExtensions)-1; i++ {
		toReturn += m.fileExtensions[i] + " "
	}
	toReturn += "or " + m.fileExtensions[i]
	return toReturn
}
