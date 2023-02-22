package app

import (
	"github.com/Zebbeni/ansizalizer/state"
	"github.com/charmbracelet/bubbles/key"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/controls"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
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
	state *state.Model

	w, h int
	km   *keyboard.Map

	controls *controls.Controls
	viewer   *viewer.Model
	help     *help.Model

	keyHandlers []keyboard.Handler

	viewport viewport.Model
}

func New() *App {
	s := state.New()

	keymap := keyboard.InitMap()

	c := controls.New(1, 1, style.ControlsBorder, keymap)
	v := viewer.New(style.ViewerBorder)

	h := help.New()
	h.ShowAll = false

	vp := viewport.New(1, 1)
	vp.Style = style.ViewportBorder

	return &App{
		state:       s,
		w:           1,
		h:           1,
		km:          keymap,
		controls:    c,
		viewer:      v,
		help:        &h,
		keyHandlers: []keyboard.Handler{c, v},
		viewport:    vp,
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
		cmd = a.HandleKeyMsg(msg)
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
	height := a.h - helpHeight

	controlsContent := a.controls.View()
	controlsWidth := lipgloss.Width(controlsContent) + 4
	controlsContent = style.ControlsBorder.Copy().Width(controlsWidth).Height(height).Render(controlsContent)

	viewerContent := a.viewer.View()
	viewerWidth := a.w - lipgloss.Width(controlsContent) - 2
	viewerContent = style.ViewerBorder.Copy().Width(viewerWidth).Height(height).Render(viewerContent)

	helpContent := a.help.View(a.km)
	helpContent = lipgloss.NewStyle().Padding(0, 0, 0, 1).Render(helpContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, controlsContent, viewerContent)
	content = lipgloss.JoinVertical(lipgloss.Top, content, helpContent)

	contentStyle := lipgloss.NewStyle().Width(a.w).Height(a.h)

	a.viewport.SetContent(contentStyle.Render(content))
	return a.viewport.View()
}

func (a *App) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, a.km.Quit):
		return tea.Quit
	}

	// check with child components here until one returns a non-nil
	// command or all keyHandlers have been hit. Need to figure out
	// how to handle events that affect the app state.
	for _, handler := range a.keyHandlers {
		if isHandled := handler.HandleKeyMsg(msg); isHandled {
			break
		}
	}

	return nil
}
