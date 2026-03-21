package process

import (
	"image"
	"image/color"
	"math"
)

// to8bit converts a 16-bit color channel value to 8-bit with proper rounding,
// avoiding the truncation error of >>8 which can be off by 1.
func to8bit(v uint32) float64 {
	return math.Round(float64(v) * 255.0 / 65535.0)
}

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
			rf, gf, bf := to8bit(r), to8bit(g), to8bit(b)

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
				A: uint8(to8bit(a)),
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
			rf := (to8bit(r) - 128) * factor + 128
			gf := (to8bit(g) - 128) * factor + 128
			bf := (to8bit(b) - 128) * factor + 128

			adjusted.Set(x, y, color.RGBA{
				R: uint8(math.Max(0, math.Min(255, rf))),
				G: uint8(math.Max(0, math.Min(255, gf))),
				B: uint8(math.Max(0, math.Min(255, bf))),
				A: uint8(to8bit(a)),
			})
		}
	}
	return adjusted
}
