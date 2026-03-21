package process

import (
	"image"
	"image/color"
	"math"
)

func adjustBrightness(img image.Image, brightness int) image.Image {
	if brightness == 0 {
		return img
	}

	bounds := img.Bounds()
	adjusted := image.NewRGBA(bounds)
	factor := float64(brightness) / 100.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf, gf, bf := float64(r>>8), float64(g>>8), float64(b>>8)

			if factor > 0 {
				rf += (255 - rf) * factor
				gf += (255 - gf) * factor
				bf += (255 - bf) * factor
			} else {
				rf += rf * factor
				gf += gf * factor
				bf += bf * factor
			}

			adjusted.Set(x, y, color.RGBA{
				R: uint8(math.Max(0, math.Min(255, rf))),
				G: uint8(math.Max(0, math.Min(255, gf))),
				B: uint8(math.Max(0, math.Min(255, bf))),
				A: uint8(a >> 8),
			})
		}
	}
	return adjusted
}

func adjustContrast(img image.Image, contrast int) image.Image {
	if contrast == 0 {
		return img
	}

	bounds := img.Bounds()
	adjusted := image.NewRGBA(bounds)
	factor := (100.0 + float64(contrast)) / 100.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf := (float64(r>>8) - 128) * factor + 128
			gf := (float64(g>>8) - 128) * factor + 128
			bf := (float64(b>>8) - 128) * factor + 128

			adjusted.Set(x, y, color.RGBA{
				R: uint8(math.Max(0, math.Min(255, rf))),
				G: uint8(math.Max(0, math.Min(255, gf))),
				B: uint8(math.Max(0, math.Min(255, bf))),
				A: uint8(a >> 8),
			})
		}
	}
	return adjusted
}
