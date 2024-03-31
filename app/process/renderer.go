package process

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
)

type Renderer struct {
	Settings             settings.Model
	shadeLightBlockFuncs map[rune]blockFunc
	shadeMedBlockFuncs   map[rune]blockFunc
	shadeHeavyBlockFuncs map[rune]blockFunc
	quarterBlockFuncs    map[rune]blockFunc
	halfBlockFuncs       map[rune]blockFunc
	fullBlockFuncs       map[rune]blockFunc
}

func New(s settings.Model) Renderer {
	m := Renderer{
		Settings: s,
	}
	m.fullBlockFuncs = m.createFullBlockFuncs()
	m.halfBlockFuncs = m.createHalfBlockFuncs()
	m.quarterBlockFuncs = m.createQuarterBlockFuncs()
	m.shadeLightBlockFuncs = m.createShadeLightFuncs()
	m.shadeMedBlockFuncs = m.createShadeMedFuncs()
	m.shadeHeavyBlockFuncs = m.createShadeHeavyFuncs()
	return m
}

type blockFunc func(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64)

func (m Renderer) createQuarterBlockFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'▀': m.calcTop,
		'▐': m.calcRight,
		'▞': m.calcDiagonal,
		'▖': m.calcBotLeft,
		'▘': m.calcTopLeft,
		'▝': m.calcTopRight,
		'▗': m.calcBotRight,
	}
}
func (m Renderer) createHalfBlockFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'▀': m.calcTop,
	}
}

func (m Renderer) createFullBlockFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'█': m.calcFull,
	}
}

func (m Renderer) createShadeLightFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'░': m.calcHeavy,
	}
}

func (m Renderer) createShadeMedFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'▒': m.calcHeavy,
	}
}

func (m Renderer) createShadeHeavyFuncs() map[rune]blockFunc {
	return map[rune]blockFunc{
		'▓': m.calcHeavy,
	}
}

func (m Renderer) getLightDarkPaletted(light, dark colorful.Color) (colorful.Color, colorful.Color) {
	_, _, p := m.Settings.Colors.GetSelected()
	colors := p.Colors()

	index := colors.Index(dark)
	paletteDark := colors.Convert(dark)

	palette := make([]color.Color, len(colors))
	copy(palette, colors)

	paletteMinusDarkest := append(colors[:index], colors[index+1:]...)
	paletteLight := paletteMinusDarkest.Convert(light)

	light, _ = colorful.MakeColor(paletteLight)
	dark, _ = colorful.MakeColor(paletteDark)

	// swap light / dark if light is darker than dark
	lightBlackDist := light.DistanceLuv(black)
	darkBlackDist := dark.DistanceLuv(black)
	if darkBlackDist > lightBlackDist {
		temp := light
		light = dark
		dark = temp
	}

	return light, dark
}

func (m Renderer) getDarkestPaletted() colorful.Color {
	if !m.Settings.Colors.IsLimited() {
		return black
	}
	_, _, p := m.Settings.Colors.GetSelected()
	colors := p.Colors()
	darkest := colors.Convert(black)
	darkestConverted, _ := colorful.MakeColor(darkest)
	return darkestConverted
}

func (m Renderer) calcLight(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	if _, _, fgBg, _ := m.Settings.Characters.Selected(); fgBg == characters.OneColor {
		avg, dist := m.avgCol(r1, r2, r3, r4)
		return avg, colorful.Color{}, math.Min(1.0, math.Abs(dist-1))
	} else {
		_, dark := lightDark(r1, r2, r3, r4)
		avg := m.avgColTrue(r1, r2, r3, r4)

		if m.Settings.Colors.IsLimited() {
			avg, dark = m.getLightDarkPaletted(avg, dark)
		}

		dist := avg.DistanceLuv(black)
		return avg, dark, math.Min(1.0, math.Abs(dist))
	}
}

func (m Renderer) calcMed(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	if _, _, fgBg, _ := m.Settings.Characters.Selected(); fgBg == characters.OneColor {
		avg, dist := m.avgCol(r1, r2, r3, r4)
		return avg, colorful.Color{}, math.Min(1.0, math.Abs(dist-1))
	} else {
		_, dark := lightDark(r1, r2, r3, r4)
		avg := m.avgColTrue(r1, r2, r3, r4)

		if m.Settings.Colors.IsLimited() {
			avg, dark = m.getLightDarkPaletted(avg, dark)
		}

		dist := avg.DistanceLuv(black)
		return avg, dark, math.Min(1.0, math.Abs(dist-0.5))
	}
}

func (m Renderer) calcHeavy(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	if _, _, fgBg, _ := m.Settings.Characters.Selected(); fgBg == characters.OneColor {
		avg, dist := m.avgCol(r1, r2, r3, r4)
		return avg, colorful.Color{}, math.Min(1.0, math.Abs(dist-1))
	} else {
		_, dark := lightDark(r1, r2, r3, r4)
		avg := m.avgColTrue(r1, r2, r3, r4)

		if m.Settings.Colors.IsLimited() {
			avg, dark = m.getLightDarkPaletted(avg, dark)
		}

		dist := avg.DistanceLuv(black)
		return avg, dark, math.Min(1.0, math.Abs(dist-1))
	}
}

func (m Renderer) calcFull(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	if _, _, fgBg, _ := m.Settings.Characters.Selected(); fgBg == characters.OneColor {
		avg, _ := m.avgCol(r1, r2, r3, r4)
		return avg, colorful.Color{}, 1.0
	} else {
		_, dark := lightDark(r1, r2, r3, r4)
		avg := m.avgColTrue(r1, r2, r3, r4)

		if m.Settings.Colors.IsLimited() {
			avg, dark = m.getLightDarkPaletted(avg, dark)
		}

		dist := avg.DistanceLuv(black)
		return avg, dark, math.Min(1.0, math.Abs(dist-1))
	}
}

func (m Renderer) calcTop(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	if r1.R == 0 && r1.G == 0 && r1.B == 0 && (r3.R != 0 || r3.G != 0 || r3.B != 0) {
		r1.R = r1.G
	}
	fg, fDist := m.avgCol(r1, r2)
	bg, bDist := m.avgCol(r3, r4)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcRight(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r2, r4)
	bg, bDist := m.avgCol(r1, r3)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcDiagonal(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r2, r3)
	bg, bDist := m.avgCol(r1, r4)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcBotLeft(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r3)
	bg, bDist := m.avgCol(r1, r2, r4)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcTopLeft(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r1)
	bg, bDist := m.avgCol(r2, r3, r4)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcTopRight(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r2)
	bg, bDist := m.avgCol(r1, r3, r4)
	return fg, bg, fDist + bDist
}

func (m Renderer) calcBotRight(r1, r2, r3, r4 colorful.Color) (colorful.Color, colorful.Color, float64) {
	fg, fDist := m.avgCol(r4)
	bg, bDist := m.avgCol(r1, r2, r3)
	return fg, bg, fDist + bDist
}
