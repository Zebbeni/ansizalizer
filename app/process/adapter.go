package process

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansiart"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

func settingsToOptions(s settings.Model) ansiart.Options {
	return settingsToOptionsWithBg(s, nil)
}

func settingsToOptionsWithBg(s settings.Model, solidBg *colorful.Color) ansiart.Options {
	isTrueColor, _, palette := s.Colors.GetSelected()
	mode, charMode, colorBg, customChars := s.Characters.Selected()
	dimType, width, height, charRatio := s.Size.Info()
	doDither, doSerpentine, ditherMatrix := s.Advanced.Dithering()

	opts := ansiart.Options{
		SizeMode:  convertSizeMode(dimType),
		Width:     width,
		Height:    height,
		CharRatio: charRatio,

		CharacterMode: convertCharMode(mode),
		AsciiCharSet:  convertAsciiCharSet(charMode),
		UnicodeCharSet: convertUnicodeCharSet(charMode),
		CustomChars:   customChars,
		ColorBg:              colorBg,
		SolidBackgroundColor: solidBg,
		SelectionMode:        convertSelectionMode(s.Characters.SelectionMode()),

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

		OutputAlpha: s.Alpha.ShouldOutputAlpha(),
		TrimAlpha:   s.Alpha.TrimAlpha(),
	}

	if !isTrueColor {
		opts.Palette = palette.Colors()
	}

	return opts
}

func convertSizeMode(m size.Mode) ansiart.SizeMode {
	switch m {
	case size.Fill:
		return ansiart.Fill
	case size.Stretch:
		return ansiart.Stretch
	default:
		return ansiart.Fit
	}
}

func convertCharMode(s characters.State) ansiart.CharacterMode {
	switch s {
	case characters.Ascii:
		return ansiart.Ascii
	case characters.Custom:
		return ansiart.Custom
	default:
		return ansiart.Unicode
	}
}

func convertAsciiCharSet(s characters.State) ansiart.AsciiCharSet {
	switch s {
	case characters.AsciiNums:
		return ansiart.AsciiNums
	case characters.AsciiSpec:
		return ansiart.AsciiSpec
	case characters.AsciiAll:
		return ansiart.AsciiAll
	default:
		return ansiart.AsciiAZ
	}
}

func convertUnicodeCharSet(s characters.State) ansiart.UnicodeCharSet {
	switch s {
	case characters.UnicodeFull:
		return ansiart.UnicodeFull
	case characters.UnicodeQuart:
		return ansiart.UnicodeQuarter
	case characters.UnicodeShadeLight:
		return ansiart.UnicodeShadeLight
	case characters.UnicodeShadeMed:
		return ansiart.UnicodeShadeMed
	case characters.UnicodeShadeHeavy:
		return ansiart.UnicodeShadeHeavy
	default:
		return ansiart.UnicodeHalf
	}
}

func convertSelectionMode(s characters.State) ansiart.SelectionMode {
	switch s {
	case characters.Sequence:
		return ansiart.Repeat
	case characters.Random:
		return ansiart.Random
	default:
		return ansiart.DarkToLight
	}
}

var samplingMap = map[resize.InterpolationFunction]ansiart.SamplingFunction{
	resize.NearestNeighbor:   ansiart.NearestNeighbor,
	resize.Bicubic:           ansiart.Bicubic,
	resize.Bilinear:          ansiart.Bilinear,
	resize.Lanczos2:          ansiart.Lanczos2,
	resize.Lanczos3:          ansiart.Lanczos3,
	resize.MitchellNetravali: ansiart.MitchellNetravali,
}

func convertSampling(f resize.InterpolationFunction) ansiart.SamplingFunction {
	if s, ok := samplingMap[f]; ok {
		return s
	}
	return ansiart.NearestNeighbor
}

type ditherEntry struct {
	matrix dither.ErrorDiffusionMatrix
	value  ansiart.DitherMatrix
}

var ditherEntries = []ditherEntry{
	{dither.FloydSteinberg, ansiart.FloydSteinberg},
	{dither.Atkinson, ansiart.Atkinson},
	{dither.Burkes, ansiart.Burkes},
	{dither.FalseFloydSteinberg, ansiart.FalseFloydSteinberg},
	{dither.JarvisJudiceNinke, ansiart.JarvisJudiceNinke},
	{dither.Sierra, ansiart.Sierra},
	{dither.Sierra2, ansiart.Sierra2},
	{dither.Sierra3, ansiart.Sierra3},
	{dither.SierraLite, ansiart.SierraLite},
	{dither.TwoRowSierra, ansiart.TwoRowSierra},
	{dither.Sierra2_4A, ansiart.Sierra2_4A},
	{dither.Simple2D, ansiart.Simple2D},
	{dither.Stucki, ansiart.Stucki},
	{dither.StevenPigeon, ansiart.StevenPigeon},
}

func convertDitherMatrix(m dither.ErrorDiffusionMatrix) ansiart.DitherMatrix {
	for _, e := range ditherEntries {
		if matrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansiart.FloydSteinberg
}

func convertDitherMode(m dithering.DitherMode) ansiart.DitherMode {
	switch m {
	case dithering.Bayer:
		return ansiart.DitherModeBayer
	case dithering.ClusteredDot:
		return ansiart.DitherModeClusteredDot
	default:
		return ansiart.DitherModeMatrix
	}
}

type clusteredDotEntry struct {
	matrix dither.OrderedDitherMatrix
	value  ansiart.ClusteredDotMatrix
}

var clusteredDotEntries = []clusteredDotEntry{
	{dither.ClusteredDot4x4, ansiart.ClusteredDot4x4},
	{dither.ClusteredDot6x6, ansiart.ClusteredDot6x6},
	{dither.ClusteredDot6x6_2, ansiart.ClusteredDot6x6_2},
	{dither.ClusteredDot6x6_3, ansiart.ClusteredDot6x6_3},
	{dither.ClusteredDot8x8, ansiart.ClusteredDot8x8},
	{dither.ClusteredDotDiagonal6x6, ansiart.ClusteredDotDiagonal6x6},
	{dither.ClusteredDotDiagonal8x8, ansiart.ClusteredDotDiagonal8x8},
	{dither.ClusteredDotDiagonal8x8_2, ansiart.ClusteredDotDiagonal8x8_2},
	{dither.ClusteredDotDiagonal8x8_3, ansiart.ClusteredDotDiagonal8x8_3},
	{dither.ClusteredDotDiagonal16x16, ansiart.ClusteredDotDiagonal16x16},
	{dither.ClusteredDotHorizontalLine, ansiart.ClusteredDotHorizontalLine},
	{dither.ClusteredDotVerticalLine, ansiart.ClusteredDotVerticalLine},
	{dither.ClusteredDotSpiral5x5, ansiart.ClusteredDotSpiral5x5},
}

func convertClusteredDotMatrix(m dither.OrderedDitherMatrix) ansiart.ClusteredDotMatrix {
	for _, e := range clusteredDotEntries {
		if orderedMatrixEqual(e.matrix, m) {
			return e.value
		}
	}
	return ansiart.ClusteredDot4x4
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
