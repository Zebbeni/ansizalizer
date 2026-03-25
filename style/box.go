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

	// Pre-pad content lines with theme bg to fill inner width
	// Manually construct the box to ensure bg covers all cells
	leftBorder := topBorderStyler(border.Left)
	rightBorder := topBorderStyler(border.Right)
	botBorderLine := botLeft + bottomBorderStyler(strings.Repeat(border.Bottom, width)) + botRight

	bgFill := BgStyle()
	contentLines := strings.Split(content, "\n")
	innerWidth := width

	var bodyLines []string
	for _, line := range contentLines {
		lineWidth := lipgloss.Width(line)
		pad := ""
		if lineWidth < innerWidth {
			pad = bgFill.Render(strings.Repeat(" ", innerWidth-lineWidth))
		}
		bodyLines = append(bodyLines, leftBorder+line+pad+rightBorder)
	}

	switch b.LabelStyle.GetAlignVertical() {
	case lipgloss.Top:
		top = topLeft + topBorderStyler(gapLeft) + renderedLabel + topBorderStyler(gapRight) + topRight
		bottom = strings.Join(bodyLines, "\n") + "\n" + botBorderLine
	case lipgloss.Bottom:
		bottom = botLeft + bottomBorderStyler(gapLeft) + renderedLabel + bottomBorderStyler(gapRight) + botRight
		top = strings.Join(bodyLines, "\n")
		// Add top border
		topBorderLine := topLeft + topBorderStyler(strings.Repeat(border.Top, width)) + topRight
		top = topBorderLine + "\n" + top
	}

	return top + "\n" + bottom
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
