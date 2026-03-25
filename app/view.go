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
	vp := viewport.New(controlsWidth, m.leftPanelHeight())

	leftContent := style.ApplyBg(m.controls.View(), controlsWidth)

	vp.SetContent(style.BgStyle().
		Width(controlsWidth).
		Height(m.leftPanelHeight()).
		Render(leftContent))
	vp.Style = style.BgStyle()
	return vp.View()
}

func (m Model) renderViewer() string {
	imgString := m.viewer.View()
	imgWidth, imgHeight := lipgloss.Size(imgString)

	imgViewer := imgString

	// only render box label border around content if big enough.
	if imgHeight > 1 && imgWidth > 4 {
		boxLabelRenderer := style.BoxWithLabel{
			BoxStyle:   style.BgStyle().BorderForeground(style.ExtraDimColor).Border(lipgloss.RoundedBorder()).BorderBackground(style.ActiveTheme.Bg),
			LabelStyle: style.BgStyle().Foreground(style.ExtraDimColor).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Bottom),
		}
		imgViewer = boxLabelRenderer.Render(fmt.Sprintf("%dx%d", imgWidth, imgHeight), imgString, imgWidth)
	}

	renderViewport := viewport.New(m.rPanelWidth()-2, m.rPanelHeight()-displayHeight-2)

	vpRightStyle := style.BgStyle().Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	rightContent := vpRightStyle.Copy().Width(m.rPanelWidth() - 2).Height(m.rPanelHeight() - 4).Render(imgViewer)
	renderViewport.SetContent(rightContent)
	renderViewport.Style = style.BgStyle()

	content := style.ApplyBg(renderViewport.View(), 0)

	return style.NormalButton.Copy().BorderForeground(style.DimmedColor1).BorderBackground(style.ActiveTheme.Bg).Render(content)
}

func (m Model) renderHelp() string {
	helpBar := help.New()
	helpBar.Styles.ShortKey = style.BgStyle().Foreground(style.DimmedColor1)
	helpBar.Styles.ShortDesc = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.ShortSeparator = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.FullKey = style.BgStyle().Foreground(style.DimmedColor1)
	helpBar.Styles.FullDesc = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.FullSeparator = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.Ellipsis = style.BgStyle().Foreground(style.ExtraDimColor)
	helpContent := style.ApplyBg(helpBar.View(event.KeyMap), 0)
	return style.BgStyle().PaddingLeft(1).Width(m.w).Render(helpContent)
}
