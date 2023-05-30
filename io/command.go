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

func DisplayCmd(toDisplay string) tea.Msg {
	return DisplayMsg(toDisplay)
}
