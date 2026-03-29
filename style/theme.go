package style

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

// ThemeName identifies a color theme.
type ThemeName int

const (
	LightOnTransparent ThemeName = iota
	LightOnDark
	DarkOnLight
	DarkOnTransparent
	LightOnDarkPaletted
	DarkOnLightPaletted
)

var ThemeNames = map[ThemeName]string{
	LightOnTransparent:  "Light on Transparent",
	LightOnDark:         "Light on Dark",
	DarkOnLight:         "Dark on Light",
	DarkOnTransparent:   "Dark on Transparent",
	LightOnDarkPaletted: "Light on Dark (Paletted)",
	DarkOnLightPaletted: "Dark on Light (Paletted)",
}

var ThemeFromName = map[string]ThemeName{
	"Light on Transparent":       LightOnTransparent,
	"Light on Dark":              LightOnDark,
	"Dark on Light":              DarkOnLight,
	"Dark on Transparent":        DarkOnTransparent,
	"Light on Dark (Paletted)":   LightOnDarkPaletted,
	"Dark on Light (Paletted)":   DarkOnLightPaletted,
}

// PaletteColors holds the current palette for paletted themes.
// Set this before calling SetTheme with a paletted theme.
var PaletteColors color.Palette

// Theme defines the complete color palette for a theme.
type Theme struct {
	Name ThemeName

	// Whether the background is transparent (no bg color applied)
	Transparent bool

	// Core colors
	Bg color.Color // app background (ignored if Transparent)
	Fg color.Color // primary foreground

	// Semantic foreground colors (all used against Bg or transparent)
	Selected  color.Color // focused/highlighted elements
	Normal    color.Color // active/confirmed elements
	Dimmed    color.Color // unfocused elements
	ExtraDim  color.Color // very subtle elements (borders, summaries)
	Subtle    color.Color // secondary text
}

var themes = map[ThemeName]Theme{
	LightOnTransparent: {
		Name:        LightOnTransparent,
		Transparent: true,
		Bg:          lipgloss.Color(""),
		Fg:          lipgloss.Color("#cccccc"),
		Selected:    lipgloss.Color("#ffffff"),
		Normal:      lipgloss.Color("#aaaaaa"),
		Dimmed:      lipgloss.Color("#777777"),
		ExtraDim:    lipgloss.Color("#444444"),
		Subtle:      lipgloss.Color("#666666"),
	},
	LightOnDark: {
		Name:        LightOnDark,
		Transparent: false,
		Bg:          lipgloss.Color("#1a1a1a"),
		Fg:          lipgloss.Color("#cccccc"),
		Selected:    lipgloss.Color("#ffffff"),
		Normal:      lipgloss.Color("#aaaaaa"),
		Dimmed:      lipgloss.Color("#777777"),
		ExtraDim:    lipgloss.Color("#444444"),
		Subtle:      lipgloss.Color("#666666"),
	},
	DarkOnLight: {
		Name:        DarkOnLight,
		Transparent: false,
		Bg:          lipgloss.Color("#e8e8e8"),
		Fg:          lipgloss.Color("#333333"),
		Selected:    lipgloss.Color("#000000"),
		Normal:      lipgloss.Color("#444444"),
		Dimmed:      lipgloss.Color("#888888"),
		ExtraDim:    lipgloss.Color("#bbbbbb"),
		Subtle:      lipgloss.Color("#999999"),
	},
	DarkOnTransparent: {
		Name:        DarkOnTransparent,
		Transparent: true,
		Bg:          lipgloss.Color(""),
		Fg:          lipgloss.Color("#333333"),
		Selected:    lipgloss.Color("#000000"),
		Normal:      lipgloss.Color("#444444"),
		Dimmed:      lipgloss.Color("#888888"),
		ExtraDim:    lipgloss.Color("#bbbbbb"),
		Subtle:      lipgloss.Color("#999999"),
	},
}

// ActiveTheme is the currently active theme.
var ActiveTheme = themes[LightOnTransparent]

// SetTheme changes the active theme and refreshes all styles.
func SetTheme(name ThemeName) {
	switch name {
	case LightOnDarkPaletted:
		ActiveTheme = buildPalettedTheme(name, true)
	case DarkOnLightPaletted:
		ActiveTheme = buildPalettedTheme(name, false)
	default:
		if t, ok := themes[name]; ok {
			ActiveTheme = t
		}
	}
	Refresh()
}

// buildPalettedTheme creates a theme from the current palette colors.
// lightOnDark=true: light fg on dark bg. lightOnDark=false: dark fg on light bg.
func buildPalettedTheme(name ThemeName, lightOnDark bool) Theme {
	if len(PaletteColors) == 0 {
		// Fall back to default non-paletted theme
		if lightOnDark {
			return themes[LightOnDark]
		}
		return themes[DarkOnLight]
	}

	darkest, lightest := findDarkestLightest(PaletteColors)

	var bgR, bgG, bgB, fgR, fgG, fgB float64
	if lightOnDark {
		bgR, bgG, bgB = colorToFloat(darkest)
		fgR, fgG, fgB = colorToFloat(lightest)
	} else {
		bgR, bgG, bgB = colorToFloat(lightest)
		fgR, fgG, fgB = colorToFloat(darkest)
	}

	return Theme{
		Name:        name,
		Transparent: false,
		Bg:          floatToHex(bgR, bgG, bgB),
		Fg:          floatToHex(fgR, fgG, fgB),
		Selected:    floatToHex(fgR, fgG, fgB),
		Normal:      lerpHex(fgR, fgG, fgB, bgR, bgG, bgB, 0.25),
		Dimmed:      lerpHex(fgR, fgG, fgB, bgR, bgG, bgB, 0.50),
		ExtraDim:    lerpHex(fgR, fgG, fgB, bgR, bgG, bgB, 0.70),
		Subtle:      lerpHex(fgR, fgG, fgB, bgR, bgG, bgB, 0.55),
	}
}

func findDarkestLightest(palette color.Palette) (color.Color, color.Color) {
	var darkest, lightest color.Color
	minLum := math.MaxFloat64
	maxLum := -1.0

	for _, c := range palette {
		r, g, b := colorToFloat(c)
		// Perceived luminance
		lum := 0.299*r + 0.587*g + 0.114*b
		if lum < minLum {
			minLum = lum
			darkest = c
		}
		if lum > maxLum {
			maxLum = lum
			lightest = c
		}
	}
	return darkest, lightest
}

func colorToFloat(c color.Color) (float64, float64, float64) {
	r, g, b, _ := c.RGBA()
	return float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0
}

func floatToHex(r, g, b float64) color.Color {
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
		int(math.Round(r*255)),
		int(math.Round(g*255)),
		int(math.Round(b*255))))
}

// lerpHex interpolates between (r1,g1,b1) and (r2,g2,b2) by t (0=first, 1=second).
func lerpHex(r1, g1, b1, r2, g2, b2, t float64) color.Color {
	return floatToHex(
		r1+(r2-r1)*t,
		g1+(g2-g1)*t,
		b1+(b2-b1)*t,
	)
}

// ColorHex returns a hex string representation of a color for comparison/logging.
func ColorHex(c color.Color) string {
	if c == nil {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

// GetTheme returns the theme for the given name.
func GetTheme(name ThemeName) Theme {
	return themes[name]
}

// ApplyBg replaces every ANSI reset in the content with a reset that
// immediately re-applies the theme background, so all whitespace
// (padding, gaps, trailing spaces) renders with the theme bg.
func ApplyBg(content string, width int) string {
	if ActiveTheme.Transparent {
		return content
	}
	rf, gf, bf := colorToFloat(ActiveTheme.Bg)
	r, g, b := int(math.Round(rf*255)), int(math.Round(gf*255)), int(math.Round(bf*255))
	bgEsc := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)

	// Replace every reset with reset+bg, so subsequent spaces carry the bg.
	// lipgloss v2 uses \x1b[m (bare SGR reset) instead of \x1b[0m.
	// Replace the longer form first to avoid double-matching.
	result := strings.ReplaceAll(content, "\x1b[0m", "\x1b[0m"+bgEsc)
	result = strings.ReplaceAll(result, "\x1b[m", "\x1b[m"+bgEsc)
	// Clean up any double bg injections from the above
	result = strings.ReplaceAll(result, bgEsc+bgEsc, bgEsc)

	// Also prepend bg to the very start so leading spaces have bg
	result = bgEsc + result

	_ = width // unused now but kept for API compatibility
	return result
}

func hexToRGB(hex string) (int, int, int) {
	if len(hex) == 7 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return 0, 0, 0
	}
	r := hexByte(hex[0:2])
	g := hexByte(hex[2:4])
	b := hexByte(hex[4:6])
	return r, g, b
}

func hexByte(s string) int {
	v := 0
	for _, c := range s {
		v *= 16
		switch {
		case c >= '0' && c <= '9':
			v += int(c - '0')
		case c >= 'a' && c <= 'f':
			v += int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			v += int(c-'A') + 10
		}
	}
	return v
}

// bg returns a style with the theme background applied (if not transparent).
func bg(s lipgloss.Style) lipgloss.Style {
	if !ActiveTheme.Transparent {
		return s.Background(ActiveTheme.Bg)
	}
	return s
}

// BgStyle returns a base style with just the theme background set.
func BgStyle() lipgloss.Style {
	s := lipgloss.NewStyle()
	if !ActiveTheme.Transparent {
		s = s.Background(ActiveTheme.Bg)
	}
	return s
}

// Refresh rebuilds all exported style variables from the active theme.
func Refresh() {
	t := ActiveTheme

	// Colors
	SelectedColor1 = t.Selected
	NormalColor1 = t.Normal
	NormalColor2 = t.Subtle
	DimmedColor1 = t.Dimmed
	DimmedColor2 = t.Subtle
	ExtraDimColor = t.ExtraDim

	// Title styles
	NormalTitle = bg(lipgloss.NewStyle().Foreground(t.Normal))
	NormalParagraph = bg(lipgloss.NewStyle().Foreground(t.Subtle))
	SelectedTitle = bg(lipgloss.NewStyle().Foreground(t.Selected))
	SelectedParagraph = bg(lipgloss.NewStyle().Foreground(t.Selected))
	DimmedTitle = bg(lipgloss.NewStyle().Foreground(t.Dimmed))
	ExtraDimTitle = bg(lipgloss.NewStyle().Foreground(t.ExtraDim))
	DimmedParagraph = bg(lipgloss.NewStyle().Foreground(t.Subtle))

	// Bordered buttons
	ActiveButton = bg(lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Normal).
		Foreground(t.Normal).
		BorderBackground(t.Bg))
	FocusButton = bg(lipgloss.NewStyle().
		BorderStyle(heavyBorder()).
		BorderForeground(t.Selected).
		Foreground(t.Selected).
		BorderBackground(t.Bg))
	NormalButton = bg(lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Dimmed).
		Foreground(t.Dimmed).
		BorderBackground(t.Bg))

	// Inline button nodes
	ActiveButtonNode = bg(lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(t.Normal))
	FocusButtonNode = bg(lipgloss.NewStyle().
		Border(heavyBorder(), false, false, false, true).
		BorderForeground(t.Selected).
		Foreground(t.Selected).
		BorderBackground(t.Bg).
		Padding(0))
	NormalButtonNode = bg(lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(t.Dimmed))

	// Tab styles
	InactiveTabStyle = bg(lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderBackground(t.Bg))
	ActiveTabStyle = bg(lipgloss.NewStyle().Border(ActiveTabBorder, true).BorderBackground(t.Bg))
	FocusTabStyle = bg(lipgloss.NewStyle().Border(FocusTabBorder, true).BorderBackground(t.Bg))
	FocusActiveTabStyle = bg(lipgloss.NewStyle().Border(FocusActiveTabBorder, true).BorderBackground(t.Bg))
	InactiveTabOnHeavyStyle = bg(lipgloss.NewStyle().Border(InactiveTabOnHeavyBorder, true).BorderBackground(t.Bg))
	TabWindowStyle = bg(lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0).BorderBackground(t.Bg))
	TabWindowHeavyStyle = bg(lipgloss.NewStyle().Align(lipgloss.Left).Border(heavyWindowBorder()).UnsetBorderTop().Padding(1, 0).BorderBackground(t.Bg))
}
