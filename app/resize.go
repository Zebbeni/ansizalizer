package app

import (
	"os"
	"time"

	"golang.org/x/term"

	tea "github.com/charmbracelet/bubbletea"
)

// There is (currently) no support on Windows for detecting resize events, so
// we instead poll at regular intervals to check if the terminal size changed.
// If a resize is detected in this way, we send a WindowSizeMsg with the new
// dimensions to bubbletea, and handle it in the Model event handler
type checkSizeMsg int

const (
	resizeCheckDuration = time.Second / 4
)

func (m Model) handleSizeMsg(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	w, h := msg.Width, msg.Height
	m.w, m.h = w, h
	m.display = m.display.SetWidth(m.rPanelWidth())

	tea.ClearScreen()
	return m, nil
}

func (m Model) handleCheckSizeMsg() (Model, tea.Cmd) {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == m.w && h == m.h {
		return m, pollForSizeChange
	}
	updateSizeCmd := func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
	return m, tea.Batch(pollForSizeChange, updateSizeCmd)
}

func pollForSizeChange() tea.Msg {
	time.Sleep(resizeCheckDuration)
	return checkSizeMsg(1)
}

func (m Model) leftPanelHeight() int {
	return m.h - helpHeight
}

func (m Model) rPanelWidth() int {
	return m.w - controlsWidth
}

func (m Model) rPanelHeight() int {
	return m.h - helpHeight
}
