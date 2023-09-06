package loader

import (
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
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

func KlarikFilmic() color.Palette {
	hexes := []string{
		"#ffffff",
		"#d6dfdf",
		"#b5c4c1",
		"#8fa6a0",
		"#6f837e",
		"#536a66",
		"#2b3b3e",
		"#162424",
		"#000000",
		"#250a1d",
		"#3f1526",
		"#5a2535",
		"#82363f",
		"#a64e54",
		"#b66868",
		"#c08780",
		"#ceaea4",
		"#b2897c",
		"#9a6a5d",
		"#7c4d3f",
		"#5b2e2b",
		"#3d181b",
		"#280b15",
		"#895938",
		"#b1834e",
		"#bb995f",
		"#caac7a",
		"#d3c59f",
		"#a8ad80",
		"#84935a",
		"#5a7645",
		"#305630",
		"#1a3725",
		"#0e2724",
		"#152f3c",
		"#2d4e59",
		"#4b7674",
		"#628e87",
		"#7ca294",
		"#a5bbae",
		"#bacbc9",
		"#a1b7bf",
		"#778faa",
		"#5e6d92",
		"#424372",
		"#352959",
		"#2c173d",
		"#492854",
		"#6e3f72",
		"#935c8d",
		"#ae7d9e",
		"#c6a7b5",
		"#ac7b90",
		"#8f516c",
		"#73415a",
		"#542846",
		"#3f1831",
	}
	return hexesToColorPalette(hexes)
}

func Mudstone() color.Palette {
	hexes := []string{
		"#1b1611",
		"#1f253c",
		"#423c32",
		"#465d32",
		"#6e3f24",
		"#6b624e",
		"#90752e",
		"#cda465",
	}
	return hexesToColorPalette(hexes)
}

func IsleOfTheDead() color.Palette {
	hexes := []string{
		"#0b0b0b",
		"#454848",
		"#4f514f",
		"#5a5a5a",
		"#666666",
		"#3e3f3f",
		"#373838",
		"#242421",
		"#2c2d25",
		"#36382a",
		"#1b1b17",
		"#313333",
		"#858585",
		"#a0a0a0",
		"#717171",
		"#2c2d2d",
		"#121210",
		"#3f4132",
		"#aeaeae",
		"#575a4a",
		"#737359",
		"#858562",
		"#93906c",
		"#686652",
		"#a9a681",
		"#48534d",
		"#252928",
		"#857d62",
		"#aea282",
		"#d0cec1",
		"#c0b9a5",
		"#58503b",
		"#7a6b54",
		"#413a28",
		"#53493a",
		"#685a44",
		"#443b2e",
		"#1a201e",
		"#362e23",
		"#7a704d",
		"#222b31",
		"#364550",
	}
	return hexesToColorPalette(hexes)
}

func hexesToColorPalette(hexes []string) color.Palette {
	var colorPalette color.Palette
	for _, h := range hexes {
		c, _ := colorful.Hex(h)
		colorPalette = append(colorPalette, c)
	}
	return colorPalette
}
