package description

import (
	"image/color"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func Palette(palette color.Palette, w, h int) string {
	runes := make([]string, len(palette)/2+1)
	rows := make([]string, h)
	for idx := 0; idx < len(palette); idx += 2 {
		var fg, bg colorful.Color
		var lipFg, lipBg lipgloss.Color

		fg, _ = colorful.MakeColor(palette[idx])
		lipFg = lipgloss.Color(fg.Hex())
		style := lipgloss.NewStyle().Foreground(lipFg)

		if idx+1 < len(palette) {
			bg, _ = colorful.MakeColor(palette[idx+1])
			lipBg = lipgloss.Color(bg.Hex())
			style = style.Copy().Background(lipBg)
		}
		runes[idx/2] = style.Render(string('▀'))
	}
	for i := 0; i < h; i++ {
		start := w * i
		stop := int(math.Min(float64(w*(i+1)), float64(len(runes))))
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, runes[start:stop]...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
