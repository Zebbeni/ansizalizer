package app

import (
	"github.com/Zebbeni/ansizalizer/component/controls"
	"github.com/Zebbeni/ansizalizer/component/viewer"
	"github.com/Zebbeni/ansizalizer/env"
	"github.com/charmbracelet/lipgloss"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/component/style"
)

const (
	ResizeCheckDuration = time.Second / 4
)

// App is the highest level component in our program. It renders a viewport that
// contains all other components. It fills the terminal window and responds to
// window resizes.
type App struct {
	w, h int
	km   KeyMap

	controls *controls.Controls
	viewer   *viewer.Viewer
	viewport viewport.Model
}

func New() *App {
	c := controls.New()
	v := viewer.New()
	vp := viewport.New(1, 1)
	vp.Style = style.Border

	return &App{
		w: 1, h: 1,
		km:       initKeymap(),
		controls: c,
		viewer:   v,
		viewport: vp,
	}
}

func (a *App) Init() tea.Cmd {
	// This initiates the polling cycle for window size updates
	// but shouldn't be necessary on non-Windows computers.
	if env.PollForSizeChange {
		return pollForSizeChange
	}
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := tea.Cmd(nil)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = a.handleSizeMsg(msg)
	case tea.KeyMsg:
		cmd = a.handleKeyMsg(msg)
	case checkSizeMsg:
		cmd = a.handleCheckSizeMsg()
	}
	return a, cmd
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
// Controls takes up a variable amount of width depending on what is displayed
// and may expand as selected menu options add submenus to the width
func (a *App) View() string {
	ctrl := a.controls.View()
	view := a.viewer.View()
	content := lipgloss.JoinHorizontal(lipgloss.Top, ctrl, " | ", view)

	contentStyle := lipgloss.NewStyle().Width(a.w).Height(a.h)

	a.viewport.SetContent(contentStyle.Render(content))
	return a.viewport.View()
}
