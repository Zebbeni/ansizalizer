package process

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansipx"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

func settingsToOptions(s settings.Model) ansipx.Options {
	return settingsToOptionsWithBg(s, nil)
}

func settingsToOptionsWithBg(s settings.Model, solidBg *colorful.Color) ansipx.Options {
	isTrueColor, _, palette := s.Colors.GetSelected()
	mode, charMode, customChars := s.Characters.Selected()
	dimType, width, height, charRatio := s.Size.Info()
	doDither, doSerpentine, ditherMatrix := s.Advanced.Dithering()

	opts := ansipx.Options{
		SizeMode:  convertSizeMode(dimType),
		Width:     width,
		Height:    height,
		CharRatio: charRatio,

		CharacterMode:        convertCharMode(mode),
		AsciiCharSet:         convertAsciiCharSet(charMode),
		UnicodeCharSet:       convertUnicodeCharSet(charMode),
		CustomChars:          customChars,
		SolidBackgroundColor: solidBg,
		SelectionMode:        convertSelectionMode(s.Characters.SelectionMode()),
		RandomSeed:           s.Characters.RandomSeed(),
		VarianceThreshold:    s.Characters.VarianceThreshold(),

		TrueColor:      isTrueColor,
		AdaptToPalette: s.Colors.AdaptToPalette(),

		Brightness: s.Adjust.Brightness(),
		Contrast:   s.Adjust.Contrast(),

		Sampling:           convertSampling(s.Advanced.SamplingFunction()),
		Dithering:          doDither,
		Serpentine:         doSerpentine,
		DitherMode:         convertDitherMode(s.Advanced.DitherMode()),
		DitherMatrix:       convertDitherMatrix(ditherMatrix),
		BayerSize:          s.Advanced.BayerSize(),
		DitherStrength:     s.Advanced.DitherStrength(),
		ClusteredDotMatrix: convertClusteredDotMatrix(s.Advanced.ClusteredDotMatrix()),

		TextStyle: ansipx.TextStyle{
			Bold:          s.TextStyle.Bold(),
			Italic:        s.TextStyle.Italic(),
			Underline:     s.TextStyle.Underline(),
			Strikethrough: s.TextStyle.Strikethrough(),
		},

		OutputAlpha:    s.Alpha.ShouldOutputAlpha(),
		TrimAlpha:      s.Alpha.TrimAlpha(),
		AlphaThreshold: s.Alpha.AlphaThreshold(),
	}

	if !isTrueColor {
		opts.Palette = palette.Colors()
	}

	return opts
}

func convertSizeMode(m size.Mode) ansipx.SizeMode {
	switch m {
	case size.Fill:
		return ansipx.Fill
	case size.Stretch:
		return ansipx.Stretch
	default:
		return ansipx.Fit
	}
}

func convertCharMode(s characters.State) ansipx.CharacterMode {
	switch s {
	case characters.Ascii:
		return ansipx.Ascii
	case characters.Custom:
		return ansipx.Custom
	default:
		return ansipx.Unicode
	}
}

func convertAsciiCharSet(s characters.State) ansipx.AsciiCharSet {
	switch s {
	case characters.AsciiNums:
		return ansipx.AsciiNums
	case characters.AsciiSpec:
		return ansipx.AsciiSpec
	case characters.AsciiAll:
		return ansipx.AsciiAll
	default:
		return ansipx.AsciiAZ
	}
}

func convertUnicodeCharSet(s characters.State) ansipx.UnicodeCharSet {
	switch s {
	case characters.UnicodeFull:
		return ansipx.UnicodeFull
	case characters.UnicodeQuart:
		return ansipx.UnicodeQuarter
	case characters.UnicodeShadeLight:
		return ansipx.UnicodeShadeLight
	case characters.UnicodeShadeMed:
		return ansipx.UnicodeShadeMed
	case characters.UnicodeShadeHeavy:
		return ansipx.UnicodeShadeHeavy
	default:
		return ansipx.UnicodeHalf
	}
}

func convertSelectionMode(s characters.State) ansipx.SelectionMode {
	switch s {
	case characters.LightVariance:
		return ansipx.LightVariance
	case characters.Sequence:
		return ansipx.Repeat
	case characters.Random:
		return ansipx.Random
	default:
		return ansipx.DarkVariance
	}
}

var samplingMap = map[resize.InterpolationFunction]ansipx.SamplingFunction{
	resize.NearestNeighbor:   ansipx.NearestNeighbor,
	resize.Bicubic:           ansipx.Bicubic,
	resize.Bilinear:          ansipx.Bilinear,
	resize.Lanczos2:          ansipx.Lanczos2,
	resize.Lanczos3:          ansipx.Lanczos3,
	resize.MitchellNetravali: ansipx.MitchellNetravali,
}

func convertSampling(f resize.InterpolationFunction) ansipx.SamplingFunction {
	if s, ok := samplingMap[f]; ok {
		return s
	}
	return ansipx.NearestNeighbor
}

type ditherEntry struct {
	matrix dither.ErrorDiffusionMatrix
	value  ansipx.DitherMatrix
}

var ditherEntries = []ditherEntry{
	{dither.FloydSteinberg, ansipx.FloydSteinberg},
	{dither.Atkinson, ansipx.Atkinson},
	{dither.Burkes, ansipx.Burkes},
	{dither.FalseFloydSteinberg, ansipx.FalseFloydSteinberg},
	{dither.JarvisJudiceNinke, ansipx.JarvisJudiceNinke},
	{dither.Sierra, ansipx.Sierra},
	{dither.Sierra2, ansipx.Sierra2},
	{dither.Sierra3, ansipx.Sierra3},
	{dither.SierraLite, ansipx.SierraLite},
	{dither.TwoRowSierra, ansipx.TwoRowSierra},
	{dither.Sierra2_4A, ansipx.Sierra2_4A},
	{dither.Simple2D, ansipx.Simple2D},
	{dither.Stucki, ansipx.Stucki},
	{dither.StevenPigeon, ansipx.StevenPigeon},
}

func convertDitherMatrix(m dither.ErrorDiffusionMatrix) ansipx.DitherMatrix {
	for _, e := range ditherEntries {
		if matrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansipx.FloydSteinberg
}

func convertDitherMode(m dithering.DitherMode) ansipx.DitherMode {
	switch m {
	case dithering.Bayer:
		return ansipx.DitherModeBayer
	case dithering.ClusteredDot:
		return ansipx.DitherModeClusteredDot
	default:
		return ansipx.DitherModeMatrix
	}
}

type clusteredDotEntry struct {
	matrix dither.OrderedDitherMatrix
	value  ansipx.ClusteredDotMatrix
}

var clusteredDotEntries = []clusteredDotEntry{
	{dither.ClusteredDot4x4, ansipx.ClusteredDot4x4},
	{dither.ClusteredDot6x6, ansipx.ClusteredDot6x6},
	{dither.ClusteredDot6x6_2, ansipx.ClusteredDot6x6_2},
	{dither.ClusteredDot6x6_3, ansipx.ClusteredDot6x6_3},
	{dither.ClusteredDot8x8, ansipx.ClusteredDot8x8},
	{dither.ClusteredDotDiagonal6x6, ansipx.ClusteredDotDiagonal6x6},
	{dither.ClusteredDotDiagonal8x8, ansipx.ClusteredDotDiagonal8x8},
	{dither.ClusteredDotDiagonal8x8_2, ansipx.ClusteredDotDiagonal8x8_2},
	{dither.ClusteredDotDiagonal8x8_3, ansipx.ClusteredDotDiagonal8x8_3},
	{dither.ClusteredDotDiagonal16x16, ansipx.ClusteredDotDiagonal16x16},
	{dither.ClusteredDotHorizontalLine, ansipx.ClusteredDotHorizontalLine},
	{dither.ClusteredDotVerticalLine, ansipx.ClusteredDotVerticalLine},
	{dither.ClusteredDotSpiral5x5, ansipx.ClusteredDotSpiral5x5},
}

func convertClusteredDotMatrix(m dither.OrderedDitherMatrix) ansipx.ClusteredDotMatrix {
	for _, e := range clusteredDotEntries {
		if orderedMatrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansipx.ClusteredDot4x4
}

func orderedMatrixEqual(a, b dither.OrderedDitherMatrix) bool {
	if len(a.Matrix) != len(b.Matrix) || a.Max != b.Max {
		return false
	}
	for i := range a.Matrix {
		if len(a.Matrix[i]) != len(b.Matrix[i]) {
			return false
		}
		for j := range a.Matrix[i] {
			if a.Matrix[i][j] != b.Matrix[i][j] {
				return false
			}
		}
	}
	return true
}

func matrixEqual(a, b dither.ErrorDiffusionMatrix) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
