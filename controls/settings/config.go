package settings

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

// Config is the JSON-serializable representation of all user-configurable settings.
type Config struct {
	Theme      string           `json:"theme,omitempty"`
	Colors     ColorsConfig     `json:"colors"`
	Characters CharactersConfig `json:"characters"`
	Size       SizeConfig       `json:"size"`
	Adjust     AdjustConfig     `json:"adjust"`
	TextStyle  TextStyleConfig  `json:"textStyle"`
	Advanced   AdvancedConfig   `json:"advanced"`
	Animation  AnimationConfig  `json:"animation"`
	Alpha      AlphaConfig      `json:"alpha"`
}

type TextStyleConfig struct {
	Bold          bool `json:"bold"`
	Italic        bool `json:"italic"`
	Underline     bool `json:"underline"`
	Strikethrough bool `json:"strikethrough"`
}

type ColorsConfig struct {
	UseTrueColor   bool     `json:"useTrueColor"`
	AdaptToPalette bool     `json:"adaptToPalette"`
	PaletteName    string   `json:"paletteName,omitempty"`
	PaletteColors []string `json:"paletteColors,omitempty"`
}

type CharactersConfig struct {
	Mode          string `json:"mode"`
	AsciiMode     string `json:"asciiMode"`
	UnicodeMode   string `json:"unicodeMode"`
	SelectionMode     string  `json:"selectionMode"`
	CustomChars       string  `json:"customChars"`
	RandomSeed        int64   `json:"randomSeed,omitempty"`
	VarianceThreshold float64 `json:"varianceThreshold,omitempty"`
}

type SizeConfig struct {
	Mode      string  `json:"mode"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	CharRatio float64 `json:"charRatio"`
}

type AdjustConfig struct {
	Brightness int `json:"brightness"`
	Contrast   int `json:"contrast"`
}

type AdvancedConfig struct {
	SamplingFunction   string `json:"samplingFunction"`
	Dithering          bool   `json:"dithering"`
	Serpentine         bool   `json:"serpentine"`
	DitherMode         string `json:"ditherMode"`
	DitherMatrix       string  `json:"ditherMatrix"`
	BayerSize          int     `json:"bayerSize"`
	DitherStrength     float32 `json:"ditherStrength"`
	ClusteredDotMatrix string  `json:"clusteredDotMatrix,omitempty"`
}

type AnimationConfig struct {
	DelayMs int `json:"delayMs"`
}

type AlphaConfig struct {
	UseAlpha       bool    `json:"useAlpha"`
	TrimAlpha      bool    `json:"trimAlpha"`
	AlphaThreshold float64 `json:"alphaThreshold,omitempty"`
}

// character state name maps for serialization
var charModeNames = map[characters.State]string{
	characters.Ascii:   "Ascii",
	characters.Unicode: "Unicode",
	characters.Custom:  "Custom",
}
var charModeFromName = map[string]characters.State{
	"Ascii":   characters.Ascii,
	"Unicode": characters.Unicode,
	"Custom":  characters.Custom,
}

var asciiModeNames = map[characters.State]string{
	characters.AsciiAz:   "AZ",
	characters.AsciiNums: "0-9",
	characters.AsciiSpec: "!$",
	characters.AsciiAll:  "All",
}
var asciiModeFromName = map[string]characters.State{
	"AZ":  characters.AsciiAz,
	"0-9": characters.AsciiNums,
	"!$":  characters.AsciiSpec,
	"All": characters.AsciiAll,
}

var unicodeModeNames = map[characters.State]string{
	characters.UnicodeFull:       "Full",
	characters.UnicodeHalf:       "Half",
	characters.UnicodeQuart:      "Quarter",
	characters.UnicodeShadeLight: "ShadeLight",
	characters.UnicodeShadeMed:   "ShadeMed",
	characters.UnicodeShadeHeavy: "ShadeHeavy",
}
var unicodeModeFromName = map[string]characters.State{
	"Full":       characters.UnicodeFull,
	"Half":       characters.UnicodeHalf,
	"Quarter":    characters.UnicodeQuart,
	"ShadeLight": characters.UnicodeShadeLight,
	"ShadeMed":   characters.UnicodeShadeMed,
	"ShadeHeavy": characters.UnicodeShadeHeavy,
}

var selectionModeNames = map[characters.State]string{
	characters.DarkVariance:  "DarkToLight",
	characters.LightVariance: "LightToDark",
	characters.Sequence:      "Sequence",
	characters.Random:        "Random",
}
var selectionModeFromName = map[string]characters.State{
	"DarkToLight": characters.DarkVariance,
	"LightToDark": characters.LightVariance,
	"Sequence":    characters.Sequence,
	"Random":      characters.Random,
}

var sizeModeNames = map[size.Mode]string{
	size.Fit:     "Fit",
	size.Fill:    "Fill",
	size.Stretch: "Stretch",
}
var sizeModeFromName = map[string]size.Mode{
	"Fit":     size.Fit,
	"Fill":    size.Fill,
	"Stretch": size.Stretch,
}

var ditherModeNames = map[dithering.DitherMode]string{
	dithering.Matrix:       "Matrix",
	dithering.Bayer:        "Bayer",
	dithering.ClusteredDot: "ClusteredDot",
}
var ditherModeFromName = map[string]dithering.DitherMode{
	"Matrix":       dithering.Matrix,
	"Bayer":        dithering.Bayer,
	"ClusteredDot": dithering.ClusteredDot,
}

func (m Model) ExportConfig() Config {
	mode, _, customChars := m.Characters.Selected()
	selMode := m.Characters.SelectionMode()
	sizeMode, width, height, charRatio := m.Size.Info()
	doDither, doSerpentine, _ := m.Advanced.Dithering()

	// Export palette colors if in palette mode
	var paletteName string
	var paletteHexColors []string
	if m.Colors.IsLimited() {
		p := m.Colors.GetCurrentPalette()
		paletteName = p.Name()
		colors := p.Colors()
		paletteHexColors = make([]string, len(colors))
		for i, c := range colors {
			cf, _ := colorful.MakeColor(c)
			paletteHexColors[i] = cf.Hex()
		}
	}

	return Config{
		Colors: ColorsConfig{
			UseTrueColor:   !m.Colors.IsLimited(),
			AdaptToPalette: m.Colors.AdaptToPalette(),
			PaletteName:   paletteName,
			PaletteColors: paletteHexColors,
		},
		Characters: CharactersConfig{
			Mode:          charModeNames[mode],
			AsciiMode:     asciiModeNames[m.Characters.AsciiMode()],
			UnicodeMode:   unicodeModeNames[m.Characters.UnicodeMode()],
			SelectionMode: selectionModeNames[selMode],
			RandomSeed:        m.Characters.RandomSeed(),
			VarianceThreshold: m.Characters.VarianceThreshold(),
			CustomChars:   string(customChars),
		},
		Size: SizeConfig{
			Mode:      sizeModeNames[sizeMode],
			Width:     width,
			Height:    height,
			CharRatio: charRatio,
		},
		Adjust: AdjustConfig{
			Brightness: m.Adjust.Brightness(),
			Contrast:   m.Adjust.Contrast(),
		},
		TextStyle: TextStyleConfig{
			Bold:          m.TextStyle.Bold(),
			Italic:        m.TextStyle.Italic(),
			Underline:     m.TextStyle.Underline(),
			Strikethrough: m.TextStyle.Strikethrough(),
		},
		Advanced: AdvancedConfig{
			SamplingFunction:   m.Advanced.SamplingFunctionName(),
			Dithering:          doDither,
			Serpentine:          doSerpentine,
			DitherMode:         ditherModeNames[m.Advanced.DitherMode()],
			DitherMatrix:       m.Advanced.DitherMatrixName(),
			BayerSize:          int(m.Advanced.BayerSize()),
			DitherStrength:     m.Advanced.DitherStrength(),
			ClusteredDotMatrix: m.Advanced.ClusteredDotName(),
		},
		Animation: AnimationConfig{
			DelayMs: m.Animation.DelayMs(),
		},
		Alpha: AlphaConfig{
			UseAlpha:       m.Alpha.ShouldOutputAlpha(),
			TrimAlpha:      m.Alpha.TrimAlpha(),
			AlphaThreshold: m.Alpha.AlphaThreshold(),
		},
	}
}

func (m *Model) ApplyConfig(cfg Config) {
	m.Colors.SetMode(cfg.Colors.UseTrueColor)
	m.Colors.SetAdaptToPalette(cfg.Colors.AdaptToPalette)
	if !cfg.Colors.UseTrueColor && len(cfg.Colors.PaletteColors) > 0 {
		p := make(color.Palette, 0, len(cfg.Colors.PaletteColors))
		for _, hex := range cfg.Colors.PaletteColors {
			c, err := colorful.Hex(hex)
			if err != nil {
				c, _ = colorful.Hex(fmt.Sprintf("#%s", hex))
			}
			p = append(p, c)
		}
		m.Colors.SetPalette(cfg.Colors.PaletteName, p)
	}

	charMode := charModeFromName[cfg.Characters.Mode]
	asciiMode, ok := asciiModeFromName[cfg.Characters.AsciiMode]
	if !ok {
		asciiMode = characters.AsciiAz
	}
	unicodeMode, ok := unicodeModeFromName[cfg.Characters.UnicodeMode]
	if !ok {
		unicodeMode = characters.UnicodeHalf
	}
	selMode := selectionModeFromName[cfg.Characters.SelectionMode]
	m.Characters.SetConfig(charMode, asciiMode, unicodeMode, selMode, cfg.Characters.CustomChars)
	if cfg.Characters.RandomSeed != 0 {
		m.Characters.SetRandomSeed(cfg.Characters.RandomSeed)
	}
	if cfg.Characters.VarianceThreshold > 0 {
		m.Characters.SetVarianceThreshold(cfg.Characters.VarianceThreshold)
	}

	sizeMode := sizeModeFromName[cfg.Size.Mode]
	m.Size.SetConfig(sizeMode, cfg.Size.Width, cfg.Size.Height, cfg.Size.CharRatio)

	m.Adjust.SetConfig(cfg.Adjust.Brightness, cfg.Adjust.Contrast)

	m.TextStyle.SetConfig(cfg.TextStyle.Bold, cfg.TextStyle.Italic, cfg.TextStyle.Underline, cfg.TextStyle.Strikethrough)

	ditherMode := ditherModeFromName[cfg.Advanced.DitherMode]
	m.Advanced.SetConfig(cfg.Advanced.SamplingFunction, cfg.Advanced.Dithering, cfg.Advanced.Serpentine, ditherMode, cfg.Advanced.DitherMatrix, cfg.Advanced.ClusteredDotMatrix, cfg.Advanced.BayerSize, cfg.Advanced.DitherStrength)

	m.Animation.SetDelayMs(cfg.Animation.DelayMs)

	m.Alpha.SetConfig(cfg.Alpha.UseAlpha, cfg.Alpha.TrimAlpha)
	if cfg.Alpha.AlphaThreshold > 0 {
		m.Alpha.SetAlphaThreshold(cfg.Alpha.AlphaThreshold)
	}

	// Theme is intentionally not saved/loaded — a bad palette + paletted theme
	// could make the UI unreadable and persist across sessions.
}

// DefaultConfig returns the default settings matching the app's initial state.
func DefaultConfig() Config {
	return Config{
		Colors: ColorsConfig{
			UseTrueColor: true,
			PaletteName:  "teal_orange",
			PaletteColors: []string{
				"#000000", "#001f2b", "#003f4f", "#005566",
				"#006666", "#338888", "#66aaaa", "#99cccc",
				"#1a0a00", "#2d1600", "#804000", "#aa5500",
				"#cc6600", "#ff9933", "#ffcc88", "#ffffff",
			},
		},
		Characters: CharactersConfig{
			Mode:          "Unicode",
			UnicodeMode:   "Half",
			AsciiMode:     "AZ",
			SelectionMode: "DarkToLight",
			CustomChars:   "/%A",
		},
		Size: SizeConfig{
			Mode:      "Fit",
			Width:     150,
			Height:    50,
			CharRatio: 0.46,
		},
		Adjust: AdjustConfig{
			Brightness: 0,
			Contrast:   0,
		},
		Advanced: AdvancedConfig{
			SamplingFunction: "Nearest Neighbor",
			Dithering:        false,
			Serpentine:        false,
			DitherMode:       "Matrix",
			DitherMatrix:     "FloydSteinberg",
			BayerSize:        4,
			DitherStrength:   1.0,
		},
		Animation: AnimationConfig{
			DelayMs: 100,
		},
		Alpha: AlphaConfig{
			UseAlpha:       true,
			TrimAlpha:      false,
			AlphaThreshold: 0.5,
		},
	}
}

func SaveConfig(cfg Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
