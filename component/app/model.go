package app

import (
	"github.com/Zebbeni/ansizalizer/state"
	"github.com/charmbracelet/bubbles/key"
	"math"
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
	ControlsWidth       = 25
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
	k := keyboard.InitMap()
	s := state.New()

	c := controls.New(s, k)
	v := viewer.New(s)

	h := help.New()
	h.ShowAll = false

	vp := viewport.New(1, 1)
	vp.Style = style.ViewportBorder

	return &App{
		state:       s,
		w:           1,
		h:           1,
		km:          k,
		controls:    c,
		viewer:      v,
		help:        &h,
		keyHandlers: []keyboard.Handler{c},
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

	controlsContentHeight := lipgloss.Height(controlsContent)
	_, yPosition := a.controls.GetActivePosition()

	activeLine := int(math.Ceil(yPosition * float64(controlsContentHeight)))
	// get the offset needed to center the active line in the viewport.
	// (subtract 2 from the height to compensate for the viewport border)
	yOffset := activeLine - ((height) / 2)
	yOffset = max(0, yOffset)

	controlsViewport := viewport.New(ControlsWidth, height)
	controlsViewport.SetContent(controlsContent)
	controlsViewport.SetYOffset(yOffset)

	controlsViewport.Style = style.ControlsBorder.Copy().Width(ControlsWidth).Height(height)
	controlsViewport.Style.GetVerticalBorderSize()

	controlsContent = controlsViewport.View()

	viewerWidth := a.w - lipgloss.Width(controlsContent)
	viewerContent := lipgloss.NewStyle().Width(viewerWidth).Height(height).Render(a.viewer.View())

	viewerViewport := viewport.New(viewerWidth, height)
	viewerViewport.SetContent(viewerContent)
	viewerViewport.Style = style.ControlsBorder.Copy().Width(viewerWidth).Height(height)
	viewerContent = viewerViewport.View()

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return b
	}
	return b
}
