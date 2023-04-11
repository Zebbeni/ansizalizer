package lospec

import (
	"fmt"
	"image/color"

	"github.com/muesli/termenv"
)

func BlackAndWhite() color.Palette {
	return color.Palette{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func AnsiVga16() color.Palette {
	return color.Palette{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 170, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 170, B: 0, A: 255},
		color.RGBA{R: 170, G: 85, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 170, A: 255},
		color.RGBA{R: 170, G: 0, B: 170, A: 255},
		color.RGBA{R: 0, G: 170, B: 170, A: 255},
		color.RGBA{R: 170, G: 170, B: 170, A: 255},
		color.RGBA{R: 85, G: 85, B: 85, A: 255},
		color.RGBA{R: 255, G: 85, B: 85, A: 255},
		color.RGBA{R: 85, G: 255, B: 85, A: 255},
		color.RGBA{R: 255, G: 255, B: 85, A: 255},
		color.RGBA{R: 85, G: 85, B: 255, A: 255},
		color.RGBA{R: 255, G: 85, B: 255, A: 255},
		color.RGBA{R: 85, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func AnsiWinConsole16() color.Palette {
	return color.Palette{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 128, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 128, B: 0, A: 255},
		color.RGBA{R: 128, G: 128, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 128, A: 255},
		color.RGBA{R: 128, G: 0, B: 128, A: 255},
		color.RGBA{R: 0, G: 128, B: 128, A: 255},
		color.RGBA{R: 192, G: 192, B: 192, A: 255},
		color.RGBA{R: 128, G: 128, B: 128, A: 255},
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 255, B: 0, A: 255},
		color.RGBA{R: 255, G: 255, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 255, A: 255},
		color.RGBA{R: 255, G: 0, B: 255, A: 255},
		color.RGBA{R: 0, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}
}

func AnsiWinPowershell16() color.Palette {
	return color.Palette{
		color.RGBA{R: 12, G: 12, B: 12, A: 255},
		color.RGBA{R: 197, G: 15, B: 31, A: 255},
		color.RGBA{R: 19, G: 161, B: 14, A: 255},
		color.RGBA{R: 193, G: 156, B: 0, A: 255},
		color.RGBA{R: 0, G: 55, B: 218, A: 255},
		color.RGBA{R: 136, G: 23, B: 152, A: 255},
		color.RGBA{R: 58, G: 150, B: 221, A: 255},
		color.RGBA{R: 204, G: 204, B: 204, A: 255},
		color.RGBA{R: 118, G: 118, B: 118, A: 255},
		color.RGBA{R: 231, G: 72, B: 86, A: 255},
		color.RGBA{R: 22, G: 198, B: 12, A: 255},
		color.RGBA{R: 249, G: 241, B: 165, A: 255},
		color.RGBA{R: 59, G: 120, B: 255, A: 255},
		color.RGBA{R: 180, G: 0, B: 158, A: 255},
		color.RGBA{R: 97, G: 214, B: 214, A: 255},
		color.RGBA{R: 242, G: 242, B: 242, A: 255},
	}
}

func Ansi16() color.Palette {
	p := make(color.Palette, 0, 16)
	for i := 0; i < 16; i++ {
		ansi := termenv.ANSI.Color(fmt.Sprintf("%d", i))
		col := termenv.ConvertToRGB(ansi)
		p = append(p, col)
	}
	return p
}

func Ansi256() color.Palette {
	p := make(color.Palette, 0, 256)
	for i := 0; i < 256; i++ {
		ansi := termenv.ANSI256.Color(fmt.Sprintf("%d", i))
		col := termenv.ConvertToRGB(ansi)
		p = append(p, col)
	}
	return p
}
