package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

const (
	displayHeight = 3
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
	imgString := m.viewer.View()
	imgWidth, imgHeight := lipgloss.Size(imgString)

	imgViewer := imgString

	// only render box label border around content if big enough.
	if imgHeight > 1 && imgWidth > 4 {
		boxLabelRenderer := style.BoxWithLabel{
			BoxStyle:   lipgloss.NewStyle().BorderForeground(style.ExtraDimColor).Border(lipgloss.RoundedBorder()),
			LabelStyle: lipgloss.NewStyle().Foreground(style.ExtraDimColor).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Bottom),
		}
		imgViewer = boxLabelRenderer.Render(fmt.Sprintf("%dx%d", imgWidth, imgHeight), imgString, imgWidth)
	}

	renderViewport := viewport.New(m.rPanelWidth()-2, m.rPanelHeight()-displayHeight-2)

	vpRightStyle := lipgloss.NewStyle().Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	rightContent := vpRightStyle.Copy().Width(m.rPanelWidth() - 2).Height(m.rPanelHeight() - 4).Render(imgViewer)
	renderViewport.SetContent(rightContent)

	content := renderViewport.View()

	return style.NormalButton.Copy().BorderForeground(style.DimmedColor1).Render(content)
}

func (m Model) renderHelp() string {
	helpBar := help.New()
	helpContent := helpBar.View(event.KeyMap)
	return lipgloss.NewStyle().PaddingLeft(1).Render(helpContent)
}
