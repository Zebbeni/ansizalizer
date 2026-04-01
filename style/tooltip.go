package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

const tooltipMaxWidth = 30

// RenderTooltip renders a tooltip box with dimmed italic text and a
// "ctrl+h to hide" footer in the bottom border.
func RenderTooltip(text string) string {
	textWidth := tooltipMaxWidth - 4 // 2 border + 2 padding
	textStyle := DimmedTitle.Copy().Italic(true).Width(textWidth)
	content := textStyle.Render(text)

	hintColor := blendColor(ExtraDimColor, ActiveTheme.Bg, 0.5)

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

	footerStyle := BgStyle().Foreground(hintColor).
		AlignHorizontal(lipgloss.Right)

	// boxWidth = textWidth + padding (2) so RenderWithFooter produces the right total
	boxWidth := textWidth + 2
	return ApplyBg(box.RenderWithFooter("", content, " (ctrl+h) ", footerStyle, boxWidth), 0)
}

// blendColor mixes two colors. t=0 returns a, t=1 returns b.
func blendColor(a, b color.Color, t float64) color.Color {
	ar, ag, ab := colorToFloat(a)
	br, bg, bb := colorToFloat(b)
	return floatToHex(
		ar+(br-ar)*t,
		ag+(bg-ag)*t,
		ab+(bb-ab)*t,
	)
}
