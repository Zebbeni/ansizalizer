package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls"
	"github.com/Zebbeni/ansizalizer/display"
	"github.com/Zebbeni/ansizalizer/env"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/viewer"
)

type State int

const (
	Main State = iota
	Browser
	Settings
)

type Model struct {
	state State

	controls controls.Model
	display  display.Model
	viewer   viewer.Model

	w, h int
}

func New() Model {
	return Model{
		state:    Main,
		controls: controls.New(controlsWidth),
		display:  display.New(),
		viewer:   viewer.New(),
		w:        100,
		h:        100,
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
	switch msg := msg.(type) {
	case checkSizeMsg:
		return m.handleCheckSizeMsg()
	case tea.WindowSizeMsg:
		return m.handleSizeMsg(msg)
	case event.StartRenderMsg:
		return m.handleStartRenderMsg()
	case event.FinishRenderMsg:
		return m.handleFinishRenderMsg(msg)
	case event.StartAdaptingMsg:
		return m.handleStartAdaptingMsg()
	case event.FinishAdaptingMsg:
		return m.handleFinishAdaptingMsg(msg)
	case event.LospecRequestMsg:
		return m.handleLospecRequestMsg(msg)
	case event.LospecResponseMsg:
		return m.handleLospecResponseMsg(msg)
	case event.DisplayMsg:
		return m.handleDisplayMsg(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Copy):
			return m.handleCopy()
		case key.Matches(msg, event.KeyMap.Save):
			return m.handleSave()
		}
	}
	return m.handleControlsUpdate(msg)
}

// View puts the whole TUI together, laid out like this:
//
//	(Left Panel)                (Right Panel)
//
// ┌────────────────┬────────────────────────────────────────┐
// │   Controls     │               Display                  │
// │                ├────────────────────────────────────────┤
// │                │               Viewer                   │
// │                │                                        │
// ├────────────────┴────────────────────────────────────────┤
// │               Help                                      │
// └─────────────────────────────────────────────────────────┘
func (m Model) View() string {
	controls := m.renderControls()
	display := m.display.View()
	viewer := m.renderViewer()
	help := m.renderHelp()

	leftPanel := controls
	rightPanel := lipgloss.JoinVertical(lipgloss.Top, display, viewer)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	all := lipgloss.JoinVertical(lipgloss.Top, panels, help)

	vp := viewport.New(m.w, m.h)
	vp.SetContent(all)
	vp.Style = lipgloss.NewStyle().Width(m.w).Height(m.h)

	return vp.View()
}
