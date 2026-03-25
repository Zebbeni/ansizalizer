package style

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ThemeName identifies a color theme.
type ThemeName int

const (
	LightOnTransparent ThemeName = iota
	LightOnDark
	DarkOnLight
	DarkOnTransparent
)

var ThemeNames = map[ThemeName]string{
	LightOnTransparent: "Light on Transparent",
	LightOnDark:        "Light on Dark",
	DarkOnLight:        "Dark on Light",
	DarkOnTransparent:  "Dark on Transparent",
}

var ThemeFromName = map[string]ThemeName{
	"Light on Transparent": LightOnTransparent,
	"Light on Dark":        LightOnDark,
	"Dark on Light":        DarkOnLight,
	"Dark on Transparent":  DarkOnTransparent,
}

// Theme defines the complete color palette for a theme.
type Theme struct {
	Name ThemeName

	// Whether the background is transparent (no bg color applied)
	Transparent bool

	// Core colors
	Bg lipgloss.Color // app background (ignored if Transparent)
	Fg lipgloss.Color // primary foreground

	// Semantic foreground colors (all used against Bg or transparent)
	Selected  lipgloss.Color // focused/highlighted elements
	Normal    lipgloss.Color // active/confirmed elements
	Dimmed    lipgloss.Color // unfocused elements
	ExtraDim  lipgloss.Color // very subtle elements (borders, summaries)
	Subtle    lipgloss.Color // secondary text
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
	if t, ok := themes[name]; ok {
		ActiveTheme = t
	}
	Refresh()
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
	hex := string(ActiveTheme.Bg)
	r, g, b := hexToRGB(hex)
	bgEsc := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)

	// Replace every reset with reset+bg, so subsequent spaces carry the bg
	result := strings.ReplaceAll(content, "\x1b[0m", "\x1b[0m"+bgEsc)

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
		BorderStyle(lipgloss.RoundedBorder()).
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
		Border(lipgloss.RoundedBorder(), false, false, false, true).
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
	TabWindowStyle = bg(lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0).BorderBackground(t.Bg))
}
