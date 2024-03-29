package dithering

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

type MatrixType int

const (
	Atkinson MatrixType = iota
	Burkes
	FloydSteinberg
	FalseFloydSteinberg
	JarvisJudiceNinke
	Sierra
	Sierra2
	Sierra3
	SierraLite
	TwoRowSierra
	Sierra2_4A
	Simple2D
	StevenPigeon
	Stucki
)

var Matrices = []MatrixType{
	Atkinson,
	Burkes,
	FloydSteinberg,
	FalseFloydSteinberg,
	JarvisJudiceNinke,
	Sierra,
	Sierra2,
	Sierra3,
	SierraLite,
	TwoRowSierra,
	Sierra2_4A,
	Simple2D,
	Stucki,
	StevenPigeon,
}

var nameMap = map[MatrixType]string{
	Atkinson:            "Atkinson",
	Burkes:              "Burkes",
	FloydSteinberg:      "FloydSteinberg",
	FalseFloydSteinberg: "FalseFloydSteinberg",
	JarvisJudiceNinke:   "JarvisJudiceNinke",
	Sierra:              "Sierra",
	Sierra2:             "Sierra2",
	Sierra3:             "Sierra3",
	SierraLite:          "SierraLite",
	TwoRowSierra:        "TwoRowSierra",
	Sierra2_4A:          "Sierra2_4A",
	Simple2D:            "Simple2D",
	Stucki:              "Stucki",
	StevenPigeon:        "StevenPigeon",
}

var errorDiffMatrixMap = map[MatrixType]dither.ErrorDiffusionMatrix{
	Atkinson:            dither.Atkinson,
	Burkes:              dither.Burkes,
	FloydSteinberg:      dither.FloydSteinberg,
	FalseFloydSteinberg: dither.FalseFloydSteinberg,
	JarvisJudiceNinke:   dither.JarvisJudiceNinke,
	Sierra:              dither.Sierra,
	Sierra2:             dither.Sierra2,
	Sierra3:             dither.Sierra3,
	SierraLite:          dither.SierraLite,
	TwoRowSierra:        dither.TwoRowSierra,
	Sierra2_4A:          dither.Sierra2_4A,
	Simple2D:            dither.Simple2D,
	Stucki:              dither.Stucki,
	StevenPigeon:        dither.StevenPigeon,
}

func newMatrixMenu(width int) list.Model {
	items := menuItems()
	return newMenu(items, width, len(items))
}

type item struct {
	Type MatrixType
}

func (i item) FilterValue() string {
	return nameMap[i.Type]
}

func (i item) Title() string {
	return nameMap[i.Type]
}

func (i item) Description() string {
	return ""
}

func menuItems() []list.Item {
	items := make([]list.Item, len(Matrices))
	for i, matrix := range Matrices {
		items[i] = item{Type: matrix}
	}
	return items
}

func newMenu(items []list.Item, width, height int) list.Model {
	l := list.New(items, NewDelegate(false), width, height/2)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowPagination(true)
	l.SetShowStatusBar(false)

	l.KeyMap.ForceQuit.Unbind()
	l.KeyMap.Quit.Unbind()

	return l
}

func NewDelegate(isActive bool) list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	if isActive {
		delegate.Styles = ItemStylesActive()
	} else {
		delegate.Styles = ItemStylesInactive()
	}
	return delegate
}

func ItemStylesActive() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 2)
	s.SelectedTitle = style.SelectedTitle.Copy().Padding(0, 1, 0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)
	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	return s
}

func ItemStylesInactive() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 2)
	s.SelectedTitle = style.NormalTitle.Copy().Padding(0, 1, 0, 2)
	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	return s
}
