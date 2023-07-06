package process

import (
	"image"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

// A list of Ascii characters by ascending brightness
var asciiChars = []rune(" `.-':_,^=;><+!rc*/z?sLTv)J7(|Fi{C}fI31tlu[neoZ5Yxjya]2ESwqkP6h9d4VpOGbUAKXHm8RD#$Bg0MNWQ%&@")
var asciiAZChars = []rune(" rczsLTvJFiCfItluneoZYxjyaESwqkPhdVpOGbUAKXHmRDBgMNWQ")
var asciiNumChars = []rune(" 7315269480")
var asciiSpecChars = []rune(" `.-':_,^=;><+!*/?)(|{}[]#$%&@")

func (m Renderer) processAscii(input image.Image) string {
	imgW, imgH := float32(input.Bounds().Dx()), float32(input.Bounds().Dy())

	dimensionType, width, height := m.Settings.Size.Info()
	if dimensionType == size.Fit {
		fitHeight := float32(width) * (imgH / imgW) * PROPORTION
		fitWidth := (float32(height) * (imgW / imgH)) / PROPORTION
		if fitHeight > float32(height) {
			width = int(fitWidth)
		} else {
			height = int(fitHeight)
		}
	}

	resizeFunc := m.Settings.Advanced.SamplingFunction()
	refImg := resize.Resize(uint(width)*2, uint(height)*2, input, resizeFunc)

	palette := m.Settings.Colors.GetCurrentPalette()

	if m.Settings.Colors.IsDithered() {
		ditherer := dither.NewDitherer(palette.Colors())
		ditherer.Matrix = m.Settings.Colors.Matrix()
		if m.Settings.Colors.IsSerpentine() {
			ditherer.Serpentine = true
		}
		refImg = ditherer.Dither(refImg)
	}

	var chars []rune
	_, charMode, useFgBg, _ := m.Settings.Characters.Selected()
	switch charMode {
	case characters.AsciiAz:
		chars = asciiAZChars
	case characters.AsciiNums:
		chars = asciiNumChars
	case characters.AsciiSpec:
		chars = asciiSpecChars
	case characters.AsciiAll:
		chars = asciiChars
	}

	content := ""
	rows := make([]string, height)
	row := make([]string, width)

	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width*2; x += 2 {
			r1, _ := colorful.MakeColor(refImg.At(x, y))
			r2, _ := colorful.MakeColor(refImg.At(x+1, y))
			r3, _ := colorful.MakeColor(refImg.At(x, y+1))
			r4, _ := colorful.MakeColor(refImg.At(x+1, y+1))

			if useFgBg == characters.TwoColor {
				fg, bg, brightness := m.fgBgBrightness(r1, r2, r3, r4)

				lipFg := lipgloss.Color(fg.Hex())
				lipBg := lipgloss.Color(bg.Hex())
				style := lipgloss.NewStyle().Foreground(lipFg).Background(lipBg)

				index := min(int(brightness*float64(len(chars))), len(chars)-1)
				char := chars[index]
				charString := string(char)

				row[x/2] = style.Render(charString)
			} else {
				fg := m.avgColTrue(r1, r2, r3, r4)
				brightness := math.Min(1.0, math.Abs(fg.DistanceLuv(black)))
				if m.Settings.Colors.IsLimited() {
					fg, _ = colorful.MakeColor(palette.Colors().Convert(fg))
				}
				lipFg := lipgloss.Color(fg.Hex())
				style := lipgloss.NewStyle().Foreground(lipFg)
				index := min(int(brightness*float64(len(chars))), len(chars)-1)
				char := chars[index]
				charString := string(char)
				row[x/2] = style.Render(charString)
			}
		}
		rows[y/2] = lipgloss.JoinHorizontal(lipgloss.Top, row...)
	}
	content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	return content
}

func (m Renderer) fgBgBrightness(c ...colorful.Color) (fg, bg colorful.Color, b float64) {
	// find the darkest and lightest among given colors
	light, dark := lightDark(c...)

	avg := m.avgColTrue(c...)
	avgCol, _ := colorful.MakeColor(avg)

	//distLight := avgCol.DistanceLuv(light)
	distDark := avgCol.DistanceLuv(dark)
	distTotal := light.DistanceLuv(dark)
	var brightness float64
	if distTotal == 0 {
		brightness = 0
	} else {
		brightness = math.Min(1.0, math.Abs(distDark/distTotal))
	}

	// if paletted:
	//   convert the darkest to its closest paletted color
	//   convert the lightest to its closest paletted color (excluding the previously found color)
	if m.Settings.Colors.IsLimited() {
		light, dark = m.getLightDarkPaletted(light, dark)
	}

	return light, dark, brightness
}

func (m Renderer) avgColTrue(colors ...colorful.Color) colorful.Color {
	rSum, gSum, bSum := 0.0, 0.0, 0.0
	for _, col := range colors {
		rSum += col.R
		gSum += col.G
		bSum += col.B
	}
	count := float64(len(colors))
	avg := colorful.Color{R: rSum / count, G: gSum / count, B: bSum / count}

	return avg
}

func lightDark(c ...colorful.Color) (light, dark colorful.Color) {
	mostLight, mostDark := 0.0, 1.0
	for _, col := range c {
		_, _, l := col.Hsl()
		if l < mostDark {
			mostDark = l
			dark = col
		}
		if l > mostLight {
			mostLight = l
			light = col
		}
	}
	return
}
