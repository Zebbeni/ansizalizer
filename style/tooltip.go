package style

import (
	"charm.land/lipgloss/v2"
)

const tooltipMaxWidth = 30

// RenderTooltip renders a tooltip box with dimmed italic text and a
// "ctrl+h to hide" footer in the bottom border.
func RenderTooltip(text string) string {
	textWidth := tooltipMaxWidth - 4 // 2 border + 2 padding
	textStyle := DimmedTitle.Italic(true).Width(textWidth).PaddingBottom(1)
	content := textStyle.Render(text)

	box := BoxWithLabel{
		BoxStyle: BgStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ExtraDimColor).
			BorderBackground(ActiveTheme.Bg).
			Padding(0, 1),
		LabelStyle: BgStyle().Foreground(ExtraDimColor).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Top),
	}

	footerStyle := BgStyle().Foreground(DimmedColor1).
		AlignHorizontal(lipgloss.Right)

	// boxWidth = textWidth + padding (2) so RenderWithFooter produces the right total
	boxWidth := textWidth + 2
	return ApplyBg(box.RenderWithFooter("", content, " (ctrl+h) ", footerStyle, boxWidth), 0)
}
