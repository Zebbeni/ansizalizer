package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	"os"
	"time"
)

// There is (currently) no support on Windows for detecting resize events, so
// we instead poll at regular intervals to check if the terminal size changed.
// If a resize is detected in this way, we send a WindowSizeMsg with the new
// dimensions to bubbletea, and handle it in the App message handler
type checkSizeMsg int

const helpHeight = 4

func (a *App) handleSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	w, h := msg.Width, msg.Height
	a.w, a.h = w, h
	a.viewport.Width, a.viewport.Height = w, h
	a.viewport.Style = a.viewport.Style.Copy().Width(w).Height(h)

	a.controls.Resize(a.w, a.h-helpHeight)

	tea.ClearScreen()
	return nil
}

func (a *App) handleCheckSizeMsg() tea.Cmd {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == a.w && h == a.h {
		return pollForSizeChange
	}
	updateSizeCmd := func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
	return tea.Batch(pollForSizeChange, updateSizeCmd)
}

func pollForSizeChange() tea.Msg {
	time.Sleep(ResizeCheckDuration)
	return checkSizeMsg(1)
}
