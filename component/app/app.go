package app

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ResizeCheckDuration = time.Second / 4
)

var (
	borderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	textStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))
)

type checkSizeMsg int

// App is the highest level component in our program. It renders a viewport
// the full width and height of the terminal containing all other components,
// and polls for size changes at regular intervals.
type App struct {
	w, h int
	km   KeyMap
	vp   viewport.Model
}

func New() *App {
	view := viewport.New(1, 1)
	view.Style = borderStyle
	m := &App{w: 1, h: 1, km: initKeymap(), vp: view}
	return m
}

func (a *App) Init() tea.Cmd {
	return pollForSizeChange
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := tea.Cmd(nil)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = a.handleKeyMsg(msg)
	case checkSizeMsg:
		cmd = a.handleCheckSizeMsg()
	}
	return a, cmd
}

func (a *App) View() string {
	text := fmt.Sprintf("%d x %d", a.w, a.h)
	textStyle = textStyle.Copy().Width(a.w).Height(a.h)
	a.vp.SetContent(textStyle.Render(text))
	return a.vp.View()
}

func (a *App) handleCheckSizeMsg() tea.Cmd {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == a.w && h == a.h {
		return pollForSizeChange
	}
	a.updateSize(w, h)
	updateSizeCmd := func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
	return tea.Batch(pollForSizeChange, updateSizeCmd)
}

func (a *App) updateSize(w, h int) {
	a.w, a.h = w, h
	a.vp.Width, a.vp.Height = w, h
	a.vp.Style = a.vp.Style.Copy().Width(w).Height(h)
	tea.ClearScreen()
}

// There is (currently) no support on Windows for detecting resize events,
// so we instead poll at regular intervals to check if the terminal size
// has changed.
func pollForSizeChange() tea.Msg {
	time.Sleep(ResizeCheckDuration)
	return checkSizeMsg(1)
}
