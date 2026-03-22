package process

import (
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansiart"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

func settingsToOptions(s settings.Model) ansiart.Options {
	isTrueColor, _, palette := s.Colors.GetSelected()
	mode, charMode, useFgBg, customChars := s.Characters.Selected()
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
		ColorMode:     convertColorMode(useFgBg),
		SelectionMode: convertSelectionMode(s.Characters.SelectionMode()),

		TrueColor: isTrueColor,

		Brightness: s.Adjust.Brightness(),
		Contrast:   s.Adjust.Contrast(),

		Sampling:     convertSampling(s.Advanced.SamplingFunction()),
		Dithering:    doDither,
		Serpentine:   doSerpentine,
		DitherMatrix: convertDitherMatrix(ditherMatrix),

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

func convertColorMode(s characters.State) ansiart.ColorMode {
	switch s {
	case characters.OneColor:
		return ansiart.OneColor
	default:
		return ansiart.TwoColor
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
