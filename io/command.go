package io

import tea "github.com/charmbracelet/bubbletea"

type StartRenderMsg bool

func StartRenderCmd() tea.Msg {
	return StartRenderMsg(true)
}

type FinishRenderMsg struct {
	FilePath  string
	ImgString string
}
