package app

import (
	"fmt"
	"github.com/Zebbeni/ansizalizer/env"
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
	vp   viewport.Model
}

func New() *App {
	view := viewport.New(1, 1)
	view.Style = style.Border
	m := &App{w: 1, h: 1, km: initKeymap(), vp: view}
	return m
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

// View draws three main components: Controls, Preview, and Help.
// Everything is rendered inside a viewport that fills the whole terminal,
// which allows us to truncate the content if the terminal windows is too small
// ┎──────────┰─────────┒
// ┃ Controls ┃ Preview ┃
// └──────────┸─────────┚
// ┃        Help        ┃
// └────────────────────┚
// Controls takes up a variable amount of width depending on what is displayed
// and may expand as selected menu options add submenus to the width
func (a *App) View() string {
	text := fmt.Sprintf("%d x %d", a.w, a.h)
	textStyle := style.Text.Copy().Width(a.w).Height(a.h)
	a.vp.SetContent(textStyle.Render(text))
	return a.vp.View()
}
