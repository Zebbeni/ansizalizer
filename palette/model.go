package palette

import (
	"image/color"
	"math"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	name   string
	colors color.Palette
	width  int
	height int
}

func New(name string, colors color.Palette, w, h int) Model {
	return Model{
		name:   name,
		colors: colors,
		width:  w,
		height: h,
	}
}

func (m Model) View() string {
	title := style.SelectedTitle.Render(m.name)
	description := m.Description()

	return lipgloss.JoinVertical(lipgloss.Top, title, description)
}

func (m Model) FilterValue() string {
	return m.name
}

func (m Model) Title() string {
	return m.name
}

func (m Model) Description() string {
	runes := make([]string, len(m.colors)/2+1)
	rows := make([]string, 0, m.height)
	for idx := 0; idx < len(m.colors); idx += 2 {
		fg, _ := colorful.MakeColor(m.colors[idx])
		blockStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fg.Hex()))

		if idx+1 < len(m.colors) {
			bg, _ := colorful.MakeColor(m.colors[idx+1])
			blockStyle = blockStyle.Copy().Background(lipgloss.Color(bg.Hex()))
		}
		runes[idx/2] = blockStyle.Render(string('▀'))
	}
	for i := 0; i < m.height; i++ {
		start := m.width * i
		if start >= len(runes) {
			break
		}
		stop := int(math.Min(float64(m.width*(i+1)), float64(len(runes))))
		rows = append(rows, "")
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, runes[start:stop]...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m Model) Name() string {
	return m.name
}

func (m Model) Colors() color.Palette {
	colorsCopy := make([]color.Color, len(m.colors))
	copy(colorsCopy, m.colors)
	return colorsCopy
}
