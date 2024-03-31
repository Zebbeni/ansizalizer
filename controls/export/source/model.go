package source

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	ExpFile State = iota
	ExpDirectory
	Input
	Browser
	SubDirsYes
	SubDirsNo
)

type Model struct {
	focus State

	doExportDirectory     bool
	includeSubdirectories bool

	Browser      browser.Model
	selectedDir  string
	selectedFile string

	ShouldClose   bool
	ShouldUnfocus bool

	IsActive bool

	width int
}

func New(w int) Model {
	browserModel := browser.New(nil, w-2)

	return Model{
		focus: ExpDirectory,

		Browser: browserModel,

		doExportDirectory:     true,
		includeSubdirectories: false,

		selectedDir:  "",
		selectedFile: "",

		width:       w,
		ShouldClose: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focus {
	case Browser:
		return m.handleSrcBrowserUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	content := make([]string, 0, 5)
	content = append(content, m.drawExportTypeOptions())

	selected := lipgloss.NewStyle().PaddingTop(1).Render(m.drawSelected())
	content = append(content, selected)

	if m.focus == Browser {
		content = append(content, m.Browser.View())
	}

	if m.doExportDirectory {
		content = append(content, m.drawSubDirOptions())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) GetSelected() (path string, isDir, useSubDirs bool) {
	if m.doExportDirectory {
		isDir = true
		path = m.selectedDir
		useSubDirs = m.includeSubdirectories
	} else {
		path = m.selectedFile
		isDir = false
		useSubDirs = false
	}
	return
}
