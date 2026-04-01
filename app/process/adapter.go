package process

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansipic"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

func settingsToOptions(s settings.Model) ansipic.Options {
	return settingsToOptionsWithBg(s, nil)
}

func settingsToOptionsWithBg(s settings.Model, solidBg *colorful.Color) ansipic.Options {
	isTrueColor, _, palette := s.Colors.GetSelected()
	mode, charMode, customChars := s.Characters.Selected()
	dimType, width, height, charRatio := s.Size.Info()
	doDither, doSerpentine, ditherMatrix := s.Advanced.Dithering()

	opts := ansipic.Options{
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

		TextStyle: ansipic.TextStyle{
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

func convertSizeMode(m size.Mode) ansipic.SizeMode {
	switch m {
	case size.Fill:
		return ansipic.Fill
	case size.Stretch:
		return ansipic.Stretch
	default:
		return ansipic.Fit
	}
}

func convertCharMode(s characters.State) ansipic.CharacterMode {
	switch s {
	case characters.Ascii:
		return ansipic.Ascii
	case characters.Custom:
		return ansipic.Custom
	default:
		return ansipic.Unicode
	}
}

func convertAsciiCharSet(s characters.State) ansipic.AsciiCharSet {
	switch s {
	case characters.AsciiNums:
		return ansipic.AsciiNums
	case characters.AsciiSpec:
		return ansipic.AsciiSpec
	case characters.AsciiAll:
		return ansipic.AsciiAll
	default:
		return ansipic.AsciiAZ
	}
}

func convertUnicodeCharSet(s characters.State) ansipic.UnicodeCharSet {
	switch s {
	case characters.UnicodeFull:
		return ansipic.UnicodeFull
	case characters.UnicodeQuart:
		return ansipic.UnicodeQuarter
	case characters.UnicodeShadeLight:
		return ansipic.UnicodeShadeLight
	case characters.UnicodeShadeMed:
		return ansipic.UnicodeShadeMed
	case characters.UnicodeShadeHeavy:
		return ansipic.UnicodeShadeHeavy
	default:
		return ansipic.UnicodeHalf
	}
}

func convertSelectionMode(s characters.State) ansipic.SelectionMode {
	switch s {
	case characters.LightVariance:
		return ansipic.LightVariance
	case characters.Sequence:
		return ansipic.Repeat
	case characters.Random:
		return ansipic.Random
	default:
		return ansipic.DarkVariance
	}
}

var samplingMap = map[resize.InterpolationFunction]ansipic.SamplingFunction{
	resize.NearestNeighbor:   ansipic.NearestNeighbor,
	resize.Bicubic:           ansipic.Bicubic,
	resize.Bilinear:          ansipic.Bilinear,
	resize.Lanczos2:          ansipic.Lanczos2,
	resize.Lanczos3:          ansipic.Lanczos3,
	resize.MitchellNetravali: ansipic.MitchellNetravali,
}

func convertSampling(f resize.InterpolationFunction) ansipic.SamplingFunction {
	if s, ok := samplingMap[f]; ok {
		return s
	}
	return ansipic.NearestNeighbor
}

type ditherEntry struct {
	matrix dither.ErrorDiffusionMatrix
	value  ansipic.DitherMatrix
}

var ditherEntries = []ditherEntry{
	{dither.FloydSteinberg, ansipic.FloydSteinberg},
	{dither.Atkinson, ansipic.Atkinson},
	{dither.Burkes, ansipic.Burkes},
	{dither.FalseFloydSteinberg, ansipic.FalseFloydSteinberg},
	{dither.JarvisJudiceNinke, ansipic.JarvisJudiceNinke},
	{dither.Sierra, ansipic.Sierra},
	{dither.Sierra2, ansipic.Sierra2},
	{dither.Sierra3, ansipic.Sierra3},
	{dither.SierraLite, ansipic.SierraLite},
	{dither.TwoRowSierra, ansipic.TwoRowSierra},
	{dither.Sierra2_4A, ansipic.Sierra2_4A},
	{dither.Simple2D, ansipic.Simple2D},
	{dither.Stucki, ansipic.Stucki},
	{dither.StevenPigeon, ansipic.StevenPigeon},
}

func convertDitherMatrix(m dither.ErrorDiffusionMatrix) ansipic.DitherMatrix {
	for _, e := range ditherEntries {
		if matrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansipic.FloydSteinberg
}

func convertDitherMode(m dithering.DitherMode) ansipic.DitherMode {
	switch m {
	case dithering.Bayer:
		return ansipic.DitherModeBayer
	case dithering.ClusteredDot:
		return ansipic.DitherModeClusteredDot
	default:
		return ansipic.DitherModeMatrix
	}
}

type clusteredDotEntry struct {
	matrix dither.OrderedDitherMatrix
	value  ansipic.ClusteredDotMatrix
}

var clusteredDotEntries = []clusteredDotEntry{
	{dither.ClusteredDot4x4, ansipic.ClusteredDot4x4},
	{dither.ClusteredDot6x6, ansipic.ClusteredDot6x6},
	{dither.ClusteredDot6x6_2, ansipic.ClusteredDot6x6_2},
	{dither.ClusteredDot6x6_3, ansipic.ClusteredDot6x6_3},
	{dither.ClusteredDot8x8, ansipic.ClusteredDot8x8},
	{dither.ClusteredDotDiagonal6x6, ansipic.ClusteredDotDiagonal6x6},
	{dither.ClusteredDotDiagonal8x8, ansipic.ClusteredDotDiagonal8x8},
	{dither.ClusteredDotDiagonal8x8_2, ansipic.ClusteredDotDiagonal8x8_2},
	{dither.ClusteredDotDiagonal8x8_3, ansipic.ClusteredDotDiagonal8x8_3},
	{dither.ClusteredDotDiagonal16x16, ansipic.ClusteredDotDiagonal16x16},
	{dither.ClusteredDotHorizontalLine, ansipic.ClusteredDotHorizontalLine},
	{dither.ClusteredDotVerticalLine, ansipic.ClusteredDotVerticalLine},
	{dither.ClusteredDotSpiral5x5, ansipic.ClusteredDotSpiral5x5},
}

func convertClusteredDotMatrix(m dither.OrderedDitherMatrix) ansipic.ClusteredDotMatrix {
	for _, e := range clusteredDotEntries {
		if orderedMatrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansipic.ClusteredDot4x4
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
