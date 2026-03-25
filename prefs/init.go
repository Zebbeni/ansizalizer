package prefs

import (
	"os"
	"path/filepath"
)

var grayscalePalette = `000000
333333
666666
999999
cccccc
ffffff
`

var tealOrangePalette = `000000
001f2b
003f4f
005566
006666
338888
66aaaa
99cccc
1a0a00
2d1600
804000
aa5500
cc6600
ff9933
ffcc88
ffffff
`

// Init creates the app data directory structure and seed files if needed.
// Returns true if this is a first-time initialization.
func Init() bool {
	settingsDir := SettingsDir()
	palettesDir := PalettesDir()

	_ = os.MkdirAll(settingsDir, 0755)
	_ = os.MkdirAll(palettesDir, 0755)

	firstRun := false

	defaultPath := filepath.Join(settingsDir, "default.json")
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		firstRun = true
	}

	if firstRun {
		writeIfMissing(filepath.Join(palettesDir, "grayscale.hex"), grayscalePalette)
		writeIfMissing(filepath.Join(palettesDir, "teal_orange.hex"), tealOrangePalette)
	}

	return firstRun
}

func writeIfMissing(path, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.WriteFile(path, []byte(content), 0644)
	}
}

// WriteDefaultSettings writes the given JSON bytes as default.json and latest.json.
func WriteDefaultSettings(data []byte) {
	settingsDir := SettingsDir()
	writeIfMissing(filepath.Join(settingsDir, "default.json"), string(data))
	writeIfMissing(filepath.Join(settingsDir, "latest.json"), string(data))
}

// LatestSettingsPath returns the path to latest.json.
func LatestSettingsPath() string {
	return filepath.Join(SettingsDir(), "latest.json")
}
