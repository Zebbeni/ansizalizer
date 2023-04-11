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
