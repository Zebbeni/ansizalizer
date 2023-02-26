package app

import (
	"github.com/Zebbeni/ansizalizer/component/browser"
	"github.com/Zebbeni/ansizalizer/component/controls"
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/Zebbeni/ansizalizer/component/viewer"
	"github.com/Zebbeni/ansizalizer/env"
	"github.com/Zebbeni/ansizalizer/io"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the highest level component in our program. It renders a viewport that
// contains all other components. It fills the terminal window and responds to
// window resizes.
type Model struct {
	// child Models
	controls controls.Model
	viewer   viewer.Model
	help     help.Model

	// shared state objects
	browserState *browser.State

	// rendering attributes
	w, h int
}

func New() Model {
	bState := browser.NewState()

	c := controls.New(&bState)
	v := viewer.New()
	h := help.New()

	h.ShowAll = false

	return Model{
		controls: c,
		viewer:   v,
		help:     h,

		browserState: &bState,
	}
}

func (m Model) Init() tea.Cmd {
	// This initiates the polling cycle for window size updates
	// but shouldn't be necessary on non-Windows computers.
	if env.PollForSizeChange {
		return pollForSizeChange
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := tea.Cmd(nil)
	switch msgType := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = m.handleSizeMsg(msgType)
	case checkSizeMsg:
		cmd = m.handleCheckSizeMsg()
	case tea.KeyMsg:
		cmd = m.HandleKeyMsg(msgType)
	}
	return m, cmd
}

func (m Model) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, io.KeyMap.Quit):
		return tea.Quit
	}

	m.controls, _ = m.controls.Update(msg)
	return nil
}

// View draws three main components: Controls, Viewer, and Help.
// Everything is rendered inside a viewport that fills the whole terminal which
// allows us to truncate the content if the terminal windows is too small.
// ┎──────────┰─────────────┒
// ┃          ┃             ┃
// ┃ Controls ┃    Viewer   ┃
// ┃          ┃             ┃
// └──────────┸─────────────┚
// ┃         Help           ┃
// └────────────────────────┚
// Model takes up a variable amount of width depending on what is displayed
// and may expand as selected menu options add submenus to the width
func (m Model) View() string {

	controlsContent := m.controls.View()
	viewerContent := m.viewer.View()

	helpContent := m.help.View(io.KeyMap)
	helpContent = lipgloss.NewStyle().Padding(0, 0, 0, 1).Render(helpContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, controlsContent, viewerContent)
	content = lipgloss.JoinVertical(lipgloss.Top, content, helpContent)

	contentStyle := lipgloss.NewStyle().Width(m.w).Height(m.h)

	vp := viewport.New(m.w, m.h)
	vp.Style = style.ViewportBorder.Copy().Width(m.w).Height(m.h)
	vp.SetContent(contentStyle.Render(content))
	return vp.View()
}
