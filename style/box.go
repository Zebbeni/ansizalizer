package style

import (
	"strings"

	"charm.land/lipgloss/v2"
)

type BoxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func NewDefaultBoxWithLabel() BoxWithLabel {
	return BoxWithLabel{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DimmedColor1),

		// You could, of course, also set background and foreground colors here
		// as well.
		LabelStyle: lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			PaddingTop(0).
			PaddingBottom(0),
	}
}

func (b BoxWithLabel) Render(label, content string, width int) string {
	var (
		// Query the box style for some of its border properties so we can
		// essentially take the top border apart and put it around the label.
		border             lipgloss.Border     = b.BoxStyle.GetBorderStyle()
		topBorderStyler    func(string ...string) string = BgStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		bottomBorderStyler func(string ...string) string = BgStyle().Foreground(b.BoxStyle.GetBorderBottomForeground()).Render
		topLeft            string              = topBorderStyler(border.TopLeft)
		topRight           string              = topBorderStyler(border.TopRight)
		botLeft            string              = bottomBorderStyler(border.BottomLeft)
		botRight           string              = bottomBorderStyler(border.BottomRight)

		renderedLabel string = b.LabelStyle.Render(label)
	)

	// Render top row with the label
	borderWidth := b.BoxStyle.GetHorizontalBorderSize()
	cellsShort := max(0, width+borderWidth-lipgloss.Width(topLeft+topRight+renderedLabel))

	gap := strings.Repeat(border.Top, cellsShort)
	var gapLeft, gapRight string
	switch b.LabelStyle.GetAlignHorizontal() {
	case lipgloss.Left:
		gapRight = gap
	case lipgloss.Right:
		gapLeft = gap
	case lipgloss.Center:
		gapLeft = strings.Repeat(border.Top, cellsShort/2)
		gapRight = strings.Repeat(border.Top, cellsShort-(cellsShort/2))
	}

	var top, bottom string

	paddedContent := content

	switch b.LabelStyle.GetAlignVertical() {
	case lipgloss.Top:
		strings.Repeat(border.Top, cellsShort)
		top = topLeft + topBorderStyler(gapLeft) + renderedLabel + topBorderStyler(gapRight) + topRight
		bottom = b.BoxStyle.
			BorderTop(false).
			Width(width + borderWidth).
			Render(paddedContent)
	case lipgloss.Bottom:
		strings.Repeat(border.Bottom, cellsShort)
		bottom = botLeft + bottomBorderStyler(gapLeft) + renderedLabel + bottomBorderStyler(gapRight) + botRight
		top = b.BoxStyle.
			BorderBottom(false).
			Width(width + borderWidth).
			Render(paddedContent)
	}

	return top + "\n" + bottom
}

// RenderWithFooter renders a box with a label in the top border and a footer
// in the bottom border. FooterStyle controls horizontal alignment of the footer.
func (b BoxWithLabel) RenderWithFooter(label, content, footer string, footerStyle lipgloss.Style, width int) string {
	var (
		border             lipgloss.Border                  = b.BoxStyle.GetBorderStyle()
		topBorderStyler    func(string ...string) string    = BgStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		bottomBorderStyler func(string ...string) string    = BgStyle().Foreground(b.BoxStyle.GetBorderBottomForeground()).Render
		topLeft            string                           = topBorderStyler(border.TopLeft)
		topRight           string                           = topBorderStyler(border.TopRight)
		botLeft            string                           = bottomBorderStyler(border.BottomLeft)
		botRight           string                           = bottomBorderStyler(border.BottomRight)
		renderedLabel      string                           = b.LabelStyle.Render(label)
		renderedFooter     string                           = footerStyle.Render(footer)
	)

	borderWidth := b.BoxStyle.GetHorizontalBorderSize()

	// Top border with label
	topShort := max(0, width+borderWidth-lipgloss.Width(topLeft+topRight+renderedLabel))
	var topGapLeft, topGapRight string
	switch b.LabelStyle.GetAlignHorizontal() {
	case lipgloss.Left:
		topGapRight = strings.Repeat(border.Top, topShort)
	case lipgloss.Right:
		topGapLeft = strings.Repeat(border.Top, topShort)
	case lipgloss.Center:
		topGapLeft = strings.Repeat(border.Top, topShort/2)
		topGapRight = strings.Repeat(border.Top, topShort-(topShort/2))
	}
	top := topLeft + topBorderStyler(topGapLeft) + renderedLabel + topBorderStyler(topGapRight) + topRight

	// Middle content (no top/bottom borders)
	middle := b.BoxStyle.
		BorderTop(false).
		BorderBottom(false).
		Width(width + borderWidth).
		Render(content)

	// Bottom border with footer
	botShort := max(0, width+borderWidth-lipgloss.Width(botLeft+botRight+renderedFooter))
	var botGapLeft, botGapRight string
	switch footerStyle.GetAlignHorizontal() {
	case lipgloss.Left:
		botGapRight = strings.Repeat(border.Bottom, botShort)
	case lipgloss.Right:
		botGapLeft = strings.Repeat(border.Bottom, botShort)
	case lipgloss.Center:
		botGapLeft = strings.Repeat(border.Bottom, botShort/2)
		botGapRight = strings.Repeat(border.Bottom, botShort-(botShort/2))
	}
	bottom := botLeft + bottomBorderStyler(botGapLeft) + renderedFooter + bottomBorderStyler(botGapRight) + botRight

	return top + "\n" + middle + "\n" + bottom
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
