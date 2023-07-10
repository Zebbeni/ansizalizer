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

func (m Renderer) processCustom(input image.Image) string {
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

	doDither, doSerpentine, matrix := m.Settings.Advanced.Dithering()
	if doDither && m.Settings.Colors.IsLimited() {
		ditherer := dither.NewDitherer(palette.Colors())
		ditherer.Matrix = matrix
		if doSerpentine {
			ditherer.Serpentine = true
		}
		refImg = ditherer.Dither(refImg)
	}

	_, _, useFgBg, chars := m.Settings.Characters.Selected()
	if len(chars) == 0 {
		return "Enter at least one custom character"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
