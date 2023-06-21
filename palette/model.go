package palette

import (
	"image/color"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
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
		var fg, bg colorful.Color
		var lipFg, lipBg lipgloss.Color

		fg, _ = colorful.MakeColor(m.colors[idx])
		lipFg = lipgloss.Color(fg.Hex())
		blockStyle := lipgloss.NewStyle().Foreground(lipFg)

		if idx+1 < len(m.colors) {
			bg, _ = colorful.MakeColor(m.colors[idx+1])
			lipBg = lipgloss.Color(bg.Hex())
			blockStyle = blockStyle.Copy().Background(lipBg)
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
	return m.colors
}
