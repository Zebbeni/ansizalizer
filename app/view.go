package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"

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

	leftContent := m.controls.View()

	vp.SetContent(style.BgStyle().
		Width(controlsWidth).
		Height(m.leftPanelHeight()).
		Render(leftContent))
	vp.Style = style.BgStyle()
	return vp.View()
}

func (m Model) renderViewer() string {
	return m.viewer.View()
}

var (
	helpCache      string
	helpCacheW     int
	helpCacheTheme style.ThemeName
)

func (m Model) renderHelp() string {
	if helpCache != "" && helpCacheW == m.w && helpCacheTheme == style.ActiveTheme.Name {
		return helpCache
	}
	helpBar := help.New()
	helpBar.Styles.ShortKey = style.BgStyle().Foreground(style.DimmedColor1)
	helpBar.Styles.ShortDesc = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.ShortSeparator = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.FullKey = style.BgStyle().Foreground(style.DimmedColor1)
	helpBar.Styles.FullDesc = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.FullSeparator = style.BgStyle().Foreground(style.ExtraDimColor)
	helpBar.Styles.Ellipsis = style.BgStyle().Foreground(style.ExtraDimColor)
	helpContent := helpBar.View(event.KeyMap)
	helpCache = style.ApplyBg(style.BgStyle().PaddingLeft(1).Width(m.w).Render(helpContent), 0)
	helpCacheW = m.w
	helpCacheTheme = style.ActiveTheme.Name
	return helpCache
}
