package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/io"
)

const (
	displayHeight = 2
	helpHeight    = 1

	controlsWidth = 30
)

func (m Model) renderControls() string {
	viewport := viewport.New(controlsWidth, m.leftPanelHeight())

	leftContent := m.controls.View()

	viewport.SetContent(lipgloss.NewStyle().
		Width(controlsWidth).
		Height(m.leftPanelHeight()).
		Render(leftContent))
	return viewport.View()
}

func (m Model) renderViewer() string {
	viewer := m.viewer.View()

	renderViewport := viewport.New(m.rPanelWidth(), m.rPanelHeight()-displayHeight)

	vpRightStyle := lipgloss.NewStyle().Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	rightContent := vpRightStyle.Copy().Width(m.rPanelWidth()).Height(m.rPanelHeight()).Render(viewer)
	renderViewport.SetContent(rightContent)
	return renderViewport.View()
}

func (m Model) renderHelp() string {
	helpBar := help.New()
	return helpBar.View(io.KeyMap)
}
