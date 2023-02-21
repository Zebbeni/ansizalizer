package app

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/controls"
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/Zebbeni/ansizalizer/component/viewer"
	"github.com/Zebbeni/ansizalizer/env"
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
	viewer   *viewer.Model
	help     help.Model

	viewport viewport.Model
}

func New() *App {
	c := controls.New(1, 1, style.ControlsBorder)
	v := viewer.New(1, 1, style.ViewerBorder)

	h := help.New()
	h.ShowAll = true

	vp := viewport.New(1, 1)
	vp.Style = style.ViewportBorder

	return &App{
		w: 1, h: 1,
		km:       initKeymap(),
		controls: c,
		viewer:   v,
		help:     h,
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

// View draws three main components: Controls, Model, and Help.
// Everything is rendered inside a viewport that fills the whole terminal which
// allows us to truncate the content if the terminal windows is too small.
// ┎──────────┰─────────────┒
// ┃          ┃             ┃
// ┃ Controls ┃    Model   ┃
// ┃          ┃             ┃
// └──────────┸─────────────┚
// ┃         Help           ┃
// └────────────────────────┚
// Controls takes up a variable amount of width depending on what is displayed
// and may expand as selected menu options add submenus to the width
func (a *App) View() string {
	ctrl := a.controls.View()
	view := a.viewer.View()

	helpText := a.help.View(a.km)
	helper := style.Help.Render(helpText)

	content := lipgloss.JoinHorizontal(lipgloss.Top, ctrl, view)
	content = lipgloss.JoinVertical(lipgloss.Top, content, helper)

	contentStyle := lipgloss.NewStyle().Width(a.w).Height(a.h)

	a.viewport.SetContent(contentStyle.Render(content))
	return a.viewport.View()
}
