package file

import (
	"image/color"
	"image/color/palette"

	"github.com/charmbracelet/bubbles/list"
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
	return ""
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
	}
}
