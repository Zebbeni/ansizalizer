package dithering

import (
	"charm.land/bubbles/v2/list"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

type MatrixType int

const (
	MatAtkinson MatrixType = iota
	MatBurkes
	MatFloydSteinberg
	MatFalseFloydSteinberg
	MatJarvisJudiceNinke
	MatSierra
	MatSierra2
	MatSierra3
	MatSierraLite
	MatTwoRowSierra
	MatSierra2_4A
	MatSimple2D
	MatStevenPigeon
	MatStucki
)

var matrixNameMap = map[MatrixType]string{
	MatAtkinson:            "Atkinson",
	MatBurkes:              "Burkes",
	MatFloydSteinberg:      "FloydSteinberg",
	MatFalseFloydSteinberg: "FalseFloydSteinberg",
	MatJarvisJudiceNinke:   "JarvisJudiceNinke",
	MatSierra:              "Sierra",
	MatSierra2:             "Sierra2",
	MatSierra3:             "Sierra3",
	MatSierraLite:          "SierraLite",
	MatTwoRowSierra:        "TwoRowSierra",
	MatSierra2_4A:          "Sierra2_4A",
	MatSimple2D:            "Simple2D",
	MatStucki:              "Stucki",
	MatStevenPigeon:        "StevenPigeon",
}

var errorDiffMatrixMap = map[MatrixType]dither.ErrorDiffusionMatrix{
	MatAtkinson:            dither.Atkinson,
	MatBurkes:              dither.Burkes,
	MatFloydSteinberg:      dither.FloydSteinberg,
	MatFalseFloydSteinberg: dither.FalseFloydSteinberg,
	MatJarvisJudiceNinke:   dither.JarvisJudiceNinke,
	MatSierra:              dither.Sierra,
	MatSierra2:             dither.Sierra2,
	MatSierra3:             dither.Sierra3,
	MatSierraLite:          dither.SierraLite,
	MatTwoRowSierra:        dither.TwoRowSierra,
	MatSierra2_4A:          dither.Sierra2_4A,
	MatSimple2D:            dither.Simple2D,
	MatStucki:              dither.Stucki,
	MatStevenPigeon:        dither.StevenPigeon,
}

type ClusteredDotType int

const (
	CDClusteredDot4x4 ClusteredDotType = iota
	CDClusteredDot6x6
	CDClusteredDot6x6_2
	CDClusteredDot6x6_3
	CDClusteredDot8x8
	CDClusteredDotDiagonal6x6
	CDClusteredDotDiagonal8x8
	CDClusteredDotDiagonal8x8_2
	CDClusteredDotDiagonal8x8_3
	CDClusteredDotDiagonal16x16
	CDClusteredDotHorizontalLine
	CDClusteredDotVerticalLine
	CDClusteredDotSpiral5x5
)

var clusteredDotNameMap = map[ClusteredDotType]string{
	CDClusteredDot4x4:            "ClusteredDot4x4",
	CDClusteredDot6x6:            "ClusteredDot6x6",
	CDClusteredDot6x6_2:          "ClusteredDot6x6_2",
	CDClusteredDot6x6_3:          "ClusteredDot6x6_3",
	CDClusteredDot8x8:            "ClusteredDot8x8",
	CDClusteredDotDiagonal6x6:    "Diagonal6x6",
	CDClusteredDotDiagonal8x8:    "Diagonal8x8",
	CDClusteredDotDiagonal8x8_2:  "Diagonal8x8_2",
	CDClusteredDotDiagonal8x8_3:  "Diagonal8x8_3",
	CDClusteredDotDiagonal16x16:  "Diagonal16x16",
	CDClusteredDotHorizontalLine: "HorizontalLine",
	CDClusteredDotVerticalLine:   "VerticalLine",
	CDClusteredDotSpiral5x5:      "Spiral5x5",
}

var clusteredDotMatrixMap = map[ClusteredDotType]dither.OrderedDitherMatrix{
	CDClusteredDot4x4:            dither.ClusteredDot4x4,
	CDClusteredDot6x6:            dither.ClusteredDot6x6,
	CDClusteredDot6x6_2:          dither.ClusteredDot6x6_2,
	CDClusteredDot6x6_3:          dither.ClusteredDot6x6_3,
	CDClusteredDot8x8:            dither.ClusteredDot8x8,
	CDClusteredDotDiagonal6x6:    dither.ClusteredDotDiagonal6x6,
	CDClusteredDotDiagonal8x8:    dither.ClusteredDotDiagonal8x8,
	CDClusteredDotDiagonal8x8_2:  dither.ClusteredDotDiagonal8x8_2,
	CDClusteredDotDiagonal8x8_3:  dither.ClusteredDotDiagonal8x8_3,
	CDClusteredDotDiagonal16x16:  dither.ClusteredDotDiagonal16x16,
	CDClusteredDotHorizontalLine: dither.ClusteredDotHorizontalLine,
	CDClusteredDotVerticalLine:   dither.ClusteredDotVerticalLine,
	CDClusteredDotSpiral5x5:      dither.ClusteredDotSpiral5x5,
}

// Matrix list items
type matrixItem struct {
	Type MatrixType
}

func (i matrixItem) FilterValue() string { return matrixNameMap[i.Type] }
func (i matrixItem) Title() string       { return matrixNameMap[i.Type] }
func (i matrixItem) Description() string { return "" }

func newMatrixMenu(width int) list.Model {
	items := make([]list.Item, 0, len(matrixNameMap))
	for _, mt := range []MatrixType{MatAtkinson, MatBurkes, MatFloydSteinberg, MatFalseFloydSteinberg, MatJarvisJudiceNinke, MatSierra, MatSierra2, MatSierra3, MatSierraLite, MatTwoRowSierra, MatSierra2_4A, MatSimple2D, MatStucki, MatStevenPigeon} {
		items = append(items, matrixItem{Type: mt})
	}
	return newMenu(items, width, len(items))
}

// Clustered dot list items
type clusteredDotItem struct {
	Type ClusteredDotType
}

func (i clusteredDotItem) FilterValue() string { return clusteredDotNameMap[i.Type] }
func (i clusteredDotItem) Title() string       { return clusteredDotNameMap[i.Type] }
func (i clusteredDotItem) Description() string { return "" }

func newClusteredDotMenu(width int) list.Model {
	items := make([]list.Item, 0, len(clusteredDotNameMap))
	for _, ct := range []ClusteredDotType{CDClusteredDot4x4, CDClusteredDot6x6, CDClusteredDot6x6_2, CDClusteredDot6x6_3, CDClusteredDot8x8, CDClusteredDotDiagonal6x6, CDClusteredDotDiagonal8x8, CDClusteredDotDiagonal8x8_2, CDClusteredDotDiagonal8x8_3, CDClusteredDotDiagonal16x16, CDClusteredDotHorizontalLine, CDClusteredDotVerticalLine, CDClusteredDotSpiral5x5} {
		items = append(items, clusteredDotItem{Type: ct})
	}
	return newMenu(items, width, len(items))
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
		Border(style.HeavyBorder(), false, false, false, true).
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
