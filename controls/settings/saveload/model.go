package saveload

import (
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Tab int

const (
	LoadTab Tab = iota
	SaveTab
)

type State int

const (
	LoadTabSelect State = iota
	SaveTabSelect
	DirInput
	DirBrowser
	FilenameForm
	SaveButton
	LoadInput
	LoadBrowser
	LoadConfirmButton
	CancelButton
)

var jsonExtensions = map[string]bool{".json": true}

type Model struct {
	tab   Tab
	focus State

	dirBrowser  browser.Model
	loadBrowser browser.Model

	selectedDir      string
	dirChosen        bool
	saveResultText   string
	selectedLoadFile string
	loadResultText   string
	filenameInput    textinput.Model

	ShouldClose bool
	IsActive    bool

	SavePath string
	LoadPath string

	statusText string

	width int
}

func New(w int, saveDir, loadDir string) Model {
	dir, _ := os.Getwd()
	if saveDir != "" {
		dir = saveDir
	}
	dirBrowser := browser.New(nil, w-2)
	if saveDir != "" {
		dirBrowser = browser.NewAtDir(saveDir, nil, w-2)
	}
	loadBrowser := browser.New(jsonExtensions, w-2)
	if loadDir != "" {
		loadBrowser = browser.NewAtDir(loadDir, jsonExtensions, w-2)
	}

	return Model{
		tab:           LoadTab,
		focus:         LoadTabSelect,
		dirBrowser:    dirBrowser,
		loadBrowser:   loadBrowser,
		selectedDir:   dir,
		filenameInput: newFilenameInput(),
		ShouldClose:   false,
		IsActive:      false,
		width:         w,
	}
}

func newFilenameInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "Filename "
	promptStyle := style.NormalButtonNode.Copy().Padding(0, 1, 0, 0)
	styles := input.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Blurred.Prompt = promptStyle
	styles.Focused.Placeholder = lipgloss.NewStyle()
	styles.Blurred.Placeholder = lipgloss.NewStyle()
	styles.Cursor.Blink = true
	input.SetStyles(styles)
	input.SetValue("settings")
	input.CharLimit = 64
	return input
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.focus {
	case DirBrowser:
		return m.handleDirBrowserUpdate(msg)
	case LoadBrowser:
		return m.handleLoadBrowserUpdate(msg)
	case FilenameForm:
		if m.filenameInput.Focused() {
			return m.handleFilenameUpdate(msg)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.drawContent()
}

func (m *Model) SetTab(loadTab bool) {
	if loadTab {
		m.tab = LoadTab
		m.focus = LoadTabSelect
	} else {
		m.tab = SaveTab
		m.focus = SaveTabSelect
	}
}

func (m *Model) ResetFocus() {
	m.focus = LoadTabSelect
	m.filenameInput.Blur()
	m.SavePath = ""
	m.LoadPath = ""
}

func (m Model) Summary() string {
	if m.statusText != "" {
		return m.statusText
	}
	return "Load/Save Settings"
}

func (m *Model) SetStatus(text string) {
	m.statusText = text
}

func (m *Model) ClearStatus() {
	m.statusText = ""
}

func (m *Model) SetLoadResult(text string) {
	m.loadResultText = text
}

func (m Model) SelectedDirPath() string {
	return m.selectedDir
}

func (m Model) SelectedLoadDir() string {
	return m.loadBrowser.SelectedDir
}

func (m *Model) SetSaveResult(text string) {
	m.saveResultText = text
	m.dirChosen = false
}
