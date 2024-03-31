package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type BoxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func NewDefaultBoxWithLabel() BoxWithLabel {
	return BoxWithLabel{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")),

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
		topBorderStyler    func(string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		bottomBorderStyler func(string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderBottomForeground()).Render
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

	switch b.LabelStyle.GetAlignVertical() {
	case lipgloss.Top:
		strings.Repeat(border.Top, cellsShort)
		top = topLeft + topBorderStyler(gapLeft) + renderedLabel + topBorderStyler(gapRight) + topRight
		bottom = b.BoxStyle.Copy().
			BorderTop(false).
			Width(width).
			Render(content)
	case lipgloss.Bottom:
		strings.Repeat(border.Bottom, cellsShort)
		bottom = botLeft + bottomBorderStyler(gapLeft) + renderedLabel + bottomBorderStyler(gapRight) + botRight
		top = b.BoxStyle.Copy().
			BorderBottom(false).
			Width(width).
			Render(content)
	}

	// Stack the pieces
	return top + "\n" + bottom
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
