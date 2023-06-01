package io

import (
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
)

type StartRenderMsg bool

func StartRenderCmd() tea.Msg {
	return StartRenderMsg(true)
}

type FinishRenderMsg struct {
	FilePath  string
	ImgString string
}

type StartAdaptingMsg bool

func StartAdaptingCmd() tea.Msg {
	return StartAdaptingMsg(true)
}

type FinishAdaptingMsg struct {
	Palette color.Palette
}

// DisplayMsg could eventually contain a type
// that indicates what style to use (warning, error, etc.)
type DisplayMsg string

func BuildDisplayCmd(msg string) tea.Cmd {
	return func() tea.Msg { return DisplayMsg(msg) }
}

func ClearDisplayCmd() tea.Msg {
	return DisplayMsg("")
}

func BuildAdaptingCmd() tea.Cmd {
	return tea.Batch(StartAdaptingCmd, BuildDisplayCmd("Generating palette..."))
}
