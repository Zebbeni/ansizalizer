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
// dimensions to bubbletea, and handle it in the Model message handler
type checkSizeMsg int

const (
	resizeCheckDuration = time.Second / 4

	helpHeight    = 2
	controlsWidth = 25
)

func (m Model) handleSizeMsg(msg tea.WindowSizeMsg) tea.Cmd {
	w, h := msg.Width, msg.Height
	m.w, m.h = w, h

	m.controls.Resize(controlsWidth, h-helpHeight)
	m.viewer.Resize(w-controlsWidth, h-helpHeight)

	tea.ClearScreen()
	return nil
}

func (m *Model) handleCheckSizeMsg() tea.Cmd {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == m.w && h == m.h {
		return pollForSizeChange
	}
	updateSizeCmd := func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
	return tea.Batch(pollForSizeChange, updateSizeCmd)
}

func pollForSizeChange() tea.Msg {
	time.Sleep(resizeCheckDuration)
	return checkSizeMsg(1)
}
