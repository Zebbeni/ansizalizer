package process

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"testing"

	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	// keep nfnt/resize for TestAvgColPaletteDirectly
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors"
)

// TestGraysPaletteOutput reproduces the user's bug: a palette with
// 000000, 626262, ffffff produces ANSI output with 97;97;97 (#616161)
// instead of 98;98;98 (#626262).
func TestGraysPaletteOutput(t *testing.T) {
	// Force true-color mode for lipgloss in test environment
	os.Setenv("COLORTERM", "truecolor")
	lipgloss.SetColorProfile(termenv.TrueColor)

	// The grays.hex palette
	black, _ := colorful.Hex("#000000")
	gray, _ := colorful.Hex("#626262")
	white, _ := colorful.Hex("#ffffff")
	paletteColors := color.Palette{black, gray, white}

	s := settings.New(40)
	s.Colors = colors.NewPaletteMode(paletteColors, 40)
	renderer := New(s)

	// Test with a solid #626262 image
	t.Run("solid_626262", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				img.Set(x, y, color.RGBA{R: 98, G: 98, B: 98, A: 255})
			}
		}

		checkOutputContainsPaletteColor(t, renderer, img, "#626262")
	})

	// Test with grays near #626262 that should map to it
	t.Run("mixed_grays_near_626262", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		grays := []uint8{90, 95, 100, 105}
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				v := grays[(y*4+x)%len(grays)]
				img.Set(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
			}
		}
		checkOutputContainsPaletteColor(t, renderer, img, "#626262")
	})

	// Test with a JPEG-like YCbCr image
	t.Run("ycbcr_626262", func(t *testing.T) {
		ycbcrImg := image.NewYCbCr(image.Rect(0, 0, 4, 4), image.YCbCrSubsampleRatio444)
		yy, cb, cr := color.RGBToYCbCr(98, 98, 98)
		for i := range ycbcrImg.Y {
			ycbcrImg.Y[i] = yy
			ycbcrImg.Cb[i] = cb
			ycbcrImg.Cr[i] = cr
		}
		checkOutputContainsPaletteColor(t, renderer, ycbcrImg, "#626262")
	})
}

func checkOutputContainsPaletteColor(t *testing.T, renderer Renderer, img image.Image, wantHex string) {
	t.Helper()

	output := renderer.process(img)

	wantColor, _ := colorful.Hex(wantHex)
	r := uint8(wantColor.R*255.0 + 0.5)
	g := uint8(wantColor.G*255.0 + 0.5)
	b := uint8(wantColor.B*255.0 + 0.5)

	// Check for the RGB values in ANSI escape sequences
	wantRGB := fmt.Sprintf("%d;%d;%d", r, g, b)

	if !strings.Contains(output, wantRGB) {
		// Find what RGB values ARE in the output
		t.Errorf("output does not contain expected RGB %s (%s)", wantRGB, wantHex)

		// Try to find nearby values
		for delta := -3; delta <= 3; delta++ {
			nearby := fmt.Sprintf("%d;%d;%d", int(r)+delta, int(g)+delta, int(b)+delta)
			if strings.Contains(output, nearby) {
				t.Errorf("  found nearby RGB: %s (off by %d)", nearby, delta)
			}
		}
	}
}

// TestAvgColPaletteDirectly tests the avgCol function directly with
// various gray inputs and the grays palette to find precision loss.
func TestAvgColPaletteDirectly(t *testing.T) {
	black, _ := colorful.Hex("#000000")
	gray, _ := colorful.Hex("#626262")
	white, _ := colorful.Hex("#ffffff")
	paletteColors := color.Palette{black, gray, white}

	s := settings.New(40)
	s.Colors = colors.NewPaletteMode(paletteColors, 40)
	renderer := New(s)

	// Test avgCol with pairs of identical gray pixels
	grayValues := []uint8{90, 95, 97, 98, 99, 100, 105}
	for _, v := range grayValues {
		t.Run(fmt.Sprintf("gray_%d", v), func(t *testing.T) {
			c := colorful.Color{R: float64(v) / 255.0, G: float64(v) / 255.0, B: float64(v) / 255.0}
			avg, _ := renderer.avgCol(c, c)
			if avg.Hex() != "#626262" {
				t.Errorf("avgCol of gray %d pair: got %s, want #626262", v, avg.Hex())
			}
		})
	}

	// Test the full getBlock path
	t.Run("getBlock_gray98", func(t *testing.T) {
		c := colorful.Color{R: 98.0 / 255.0, G: 98.0 / 255.0, B: 98.0 / 255.0}
		_, fg, bg := renderer.getBlock(c, c, c, c)

		// Extra MakeColor round-trip as processUnicode does
		pFg, _ := colorful.MakeColor(fg)
		pBg, _ := colorful.MakeColor(bg)

		if pFg.Hex() != "#626262" {
			t.Errorf("getBlock fg: got %s, want #626262", pFg.Hex())
		}
		if pBg.Hex() != "#626262" {
			t.Errorf("getBlock bg: got %s, want #626262", pBg.Hex())
		}
	})

	// Test with pixels from a resized image
	t.Run("resized_pixels", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				img.Set(x, y, color.RGBA{R: 98, G: 98, B: 98, A: 255})
			}
		}
		resizeFunc := renderer.Settings.Advanced.SamplingFunction()
		refImg := resize.Resize(4, 4, img, resizeFunc)

		r1, _ := colorful.MakeColor(refImg.At(0, 0))
		r2, _ := colorful.MakeColor(refImg.At(1, 0))

		t.Logf("resize output type: %T", refImg)
		t.Logf("r1: R=%v G=%v B=%v hex=%s", r1.R, r1.G, r1.B, r1.Hex())

		avg, _ := renderer.avgCol(r1, r2)
		t.Logf("avgCol result: hex=%s R=%v", avg.Hex(), avg.R)

		pAvg, _ := colorful.MakeColor(avg)
		t.Logf("after MakeColor: hex=%s R=%v", pAvg.Hex(), pAvg.R)
	})
}

// TestAdjust16BitPrecision verifies that adjustBrightness and adjustContrast
// properly round 16-bit pixel values to 8-bit.
func TestAdjust16BitPrecision(t *testing.T) {
	img := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	img.SetNRGBA64(0, 0, color.NRGBA64{R: 0x62FF, G: 0x62FF, B: 0x62FF, A: 0xFFFF})

	result := adjustBrightness(img, 1)
	col, _ := colorful.MakeColor(result.At(0, 0))
	if col.Hex() != "#646464" {
		t.Errorf("brightness=1 on 16-bit 0x62FF: got %s, want #646464", col.Hex())
	}

	result = adjustContrast(img, 1)
	col, _ = colorful.MakeColor(result.At(0, 0))
	if col.Hex() != "#626262" {
		t.Errorf("contrast=1 on 16-bit 0x62FF: got %s, want #626262", col.Hex())
	}
}
