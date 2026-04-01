package style

import (
	"testing"
)

func TestSetThemeLightOnDark(t *testing.T) {
	SetTheme(LightOnDark)
	if ActiveTheme.Name != LightOnDark {
		t.Errorf("ActiveTheme.Name = %v, want LightOnDark", ActiveTheme.Name)
	}
	if ActiveTheme.Transparent {
		t.Error("LightOnDark should not be transparent")
	}
}

func TestSetThemeLightOnTransparent(t *testing.T) {
	SetTheme(LightOnTransparent)
	if ActiveTheme.Name != LightOnTransparent {
		t.Errorf("ActiveTheme.Name = %v, want LightOnTransparent", ActiveTheme.Name)
	}
	if !ActiveTheme.Transparent {
		t.Error("LightOnTransparent should be transparent")
	}
}

func TestSetThemeDarkOnLight(t *testing.T) {
	SetTheme(DarkOnLight)
	if ActiveTheme.Name != DarkOnLight {
		t.Errorf("ActiveTheme.Name = %v, want DarkOnLight", ActiveTheme.Name)
	}
}

func TestColorHexNil(t *testing.T) {
	if got := ColorHex(nil); got != "" {
		t.Errorf("ColorHex(nil) = %q, want empty", got)
	}
}

func TestApplyBgTransparentIsNoop(t *testing.T) {
	SetTheme(LightOnTransparent)
	input := "hello \x1b[0m world"
	if got := ApplyBg(input, 0); got != input {
		t.Errorf("ApplyBg on transparent theme should return input unchanged")
	}
}

func TestApplyBgInjectsBackground(t *testing.T) {
	SetTheme(LightOnDark)
	input := "hello\x1b[0mworld"
	result := ApplyBg(input, 0)
	if result == input {
		t.Error("ApplyBg on non-transparent theme should modify the string")
	}
	// Should contain a 48;2 (bg color) escape after the reset
	if len(result) <= len(input) {
		t.Error("ApplyBg should inject background escape sequences")
	}
}

func TestThemeFromNameMapsAllThemes(t *testing.T) {
	for name, themeName := range ThemeFromName {
		if displayName, ok := ThemeNames[themeName]; !ok {
			t.Errorf("ThemeNames missing entry for %v (from ThemeFromName[%q])", themeName, name)
		} else if displayName != name {
			t.Errorf("ThemeNames[%v] = %q, want %q", themeName, displayName, name)
		}
	}
}
