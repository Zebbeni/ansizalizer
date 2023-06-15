package loader

import (
	"image/color"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/browser"
)

var (
	paletteExtensions = []string{".hex"}
)

type Model struct {
	FileBrowser browser.Model

	name    string
	palette color.Palette

	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	fileBrowser := browser.New(paletteExtensions, w-2)

	return Model{
		FileBrowser:   fileBrowser,
		ShouldUnfocus: false,
		width:         w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
		m.ShouldUnfocus = true
	}
	return m, cmd
}

func (m Model) View() string {
	return m.FileBrowser.View()
}

func (m Model) GetCurrent() color.Palette {
	return m.palette
}
