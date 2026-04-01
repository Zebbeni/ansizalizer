package palette

import (
	"image/color"
	"math"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
)

func Palette(palette color.Palette, w, h int) string {
	runes := make([]string, len(palette)/2+1)
	rows := make([]string, 0, h)
	for idx := 0; idx < len(palette); idx += 2 {
		fg, _ := colorful.MakeColor(palette[idx])
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(fg.Hex()))

		if idx+1 < len(palette) {
			bg, _ := colorful.MakeColor(palette[idx+1])
			style = style.Background(lipgloss.Color(bg.Hex()))
		}
		runes[idx/2] = style.Render(string('▀'))
	}
	for i := 0; i < h; i++ {
		start := w * i
		if start >= len(runes) {
			break
		}
		stop := int(math.Min(float64(w*(i+1)), float64(len(runes))))
		rows = append(rows, "")
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, runes[start:stop]...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
