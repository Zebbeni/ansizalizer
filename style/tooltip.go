package style

import (
	"charm.land/lipgloss/v2"
)

const tooltipMaxWidth = 30

// RenderTooltip renders a tooltip box with dimmed italic text.
func RenderTooltip(text string) string {
	textStyle := DimmedTitle.Copy().Italic(true).Width(tooltipMaxWidth - 4)
	content := textStyle.Render(text)

	border := BgStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ExtraDimColor).
		BorderBackground(ActiveTheme.Bg).
		Padding(0, 1)

	return ApplyBg(border.Render(content), 0)
}
