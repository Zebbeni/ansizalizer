package saveload

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleDirBrowserUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.dirBrowser, cmd = m.dirBrowser.Update(msg)
	m.selectedDir = m.dirBrowser.SelectedDir

	if m.dirBrowser.ShouldClose {
		m.dirChosen = true
		m.focus = FilenameForm
		m.dirBrowser.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleLoadBrowserUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.loadBrowser, cmd = m.loadBrowser.Update(msg)

	if m.loadBrowser.ShouldClose {
		if m.loadBrowser.SelectedFile != "" {
			m.selectedLoadFile = m.loadBrowser.SelectedFile
			m.focus = LoadConfirmButton
		} else {
			m.focus = LoadInput
		}
		m.loadBrowser.ShouldClose = false
	}
	return m, cmd
}

func (m Model) handleFilenameUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.filenameInput.Blur()
			m.focus = SaveButton
			return m, nil
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.filenameInput.Blur()
			m.focus = FilenameForm
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.filenameInput, cmd = m.filenameInput.Update(msg)
	return m, cmd
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case SaveTabSelect:
		m.tab = SaveTab
		m.focus = DirInput
	case LoadTabSelect:
		m.tab = LoadTab
		m.focus = LoadInput
	case DirInput:
		m.focus = DirBrowser
		m.saveResultText = ""
		m.dirBrowser = browser.NewAtDir(m.selectedDir, nil, m.width-2)
	case LoadInput:
		m.focus = LoadBrowser
		m.loadResultText = ""
		m.loadBrowser = browser.NewAtDir(m.loadBrowser.SelectedDir, jsonExtensions, m.width-2)
	case FilenameForm:
		m.filenameInput.Focus()
	case SaveButton:
		filename := m.filenameInput.Value()
		if filename == "" {
			filename = "settings"
		}
		m.SavePath = filepath.Join(m.selectedDir, filename+".json")
		m.saveResultText = ""
		m.ShouldClose = true
	case LoadConfirmButton:
		m.LoadPath = m.selectedLoadFile
		m.selectedLoadFile = ""
		m.loadResultText = ""
		m.ShouldClose = true
	case CancelButton:
		if m.tab == LoadTab {
			m.selectedLoadFile = ""
			m.focus = LoadInput
		} else {
			m.dirChosen = false
			m.focus = DirInput
		}
	}
	return m, nil
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Left):
		switch m.focus {
		case SaveTabSelect:
			m.focus = LoadTabSelect
			m.tab = LoadTab
		case CancelButton:
			if m.tab == SaveTab {
				m.focus = SaveButton
			} else if m.selectedLoadFile != "" {
				m.focus = LoadConfirmButton
			}
		case LoadConfirmButton:
			// no-op, already leftmost in confirm row
		}
	case key.Matches(msg, event.KeyMap.Right):
		switch m.focus {
		case LoadTabSelect:
			m.focus = SaveTabSelect
			m.tab = SaveTab
		case SaveButton:
			m.focus = CancelButton
		case LoadConfirmButton:
			m.focus = CancelButton
		}
	case key.Matches(msg, event.KeyMap.Down):
		switch m.focus {
		case LoadTabSelect:
			m.focus = LoadInput
		case SaveTabSelect:
			m.focus = DirInput
		case DirInput:
			if m.dirChosen {
				m.focus = FilenameForm
			}
		case FilenameForm:
			if m.filenameInput.Value() != "" {
				m.focus = SaveButton
			}
		case SaveButton, CancelButton, LoadConfirmButton:
			m.ShouldClose = true
		case LoadInput:
			if m.selectedLoadFile != "" {
				m.focus = LoadConfirmButton
			}
		}
	case key.Matches(msg, event.KeyMap.Up):
		switch m.focus {
		case LoadInput:
			m.focus = LoadTabSelect
		case DirInput:
			m.focus = SaveTabSelect
		case FilenameForm:
			m.focus = DirInput
		case SaveButton, CancelButton:
			if m.tab == SaveTab {
				m.focus = FilenameForm
			} else {
				m.focus = LoadInput
			}
		case LoadConfirmButton:
			m.focus = LoadInput
		case LoadTabSelect, SaveTabSelect:
			m.ShouldClose = true
		}
	}
	return m, nil
}
