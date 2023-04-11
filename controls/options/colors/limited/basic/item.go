package basic

import (
	"image/color"
	"image/color/palette"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

type item struct {
	name    string
	palette color.Palette
}

func (i item) FilterValue() string {
	return i.name
}

func (i item) Title() string {
	return i.name
}

func (i item) Description() string {
	blocks := make([]string, len(i.palette)/2+1)
	for idx := 0; idx < len(i.palette); idx += 2 {
		var fg, bg colorful.Color
		var lipFg, lipBg lipgloss.Color

		fg, _ = colorful.MakeColor(i.palette[idx])
		lipFg = lipgloss.Color(fg.Hex())
		style := lipgloss.NewStyle().Foreground(lipFg)

		if idx+1 < len(i.palette) {
			bg, _ = colorful.MakeColor(i.palette[idx+1])
			lipBg = lipgloss.Color(bg.Hex())
			style = style.Copy().Background(lipBg)
		}

		blocks[idx/2] = style.Render(string('▀'))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, blocks...)
}

func menuItems() []list.Item {
	return []list.Item{
		item{name: "WebSafe", palette: palette.WebSafe},
		item{name: "Plan9", palette: palette.Plan9},
		item{name: "AnsiWinConsole16", palette: AnsiWinConsole16()},
		item{name: "AnsiVGA16", palette: AnsiVga16()},
		item{name: "Ansi16", palette: Ansi16()},
		item{name: "Ansi256", palette: Ansi256()},
		item{name: "BlackAndWhite", palette: BlackAndWhite()},
		item{name: "KlarikFilmic", palette: KlarikFilmic()},
		item{name: "Mudstone", palette: Mudstone()},
		item{name: "IsleOfTheDead", palette: IsleOfTheDead()},
	}
}
