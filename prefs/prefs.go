package prefs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const prefsFileName = "prefs.json"

type Dirs struct {
	Browse       string `json:"browse,omitempty"`
	ExportSource string `json:"exportSource,omitempty"`
	ExportDest   string `json:"exportDest,omitempty"`
	SaveDir      string `json:"saveDir,omitempty"`
	LoadFile     string `json:"loadFile,omitempty"`
	PaletteLoad  string `json:"paletteLoad,omitempty"`
}

// AppPrefs stores user preferences that persist across sessions.
type AppPrefs struct {
	ShowHelp        bool   `json:"showHelp"`
	RestoreTheme    bool   `json:"restoreTheme"`
	RestoreSettings bool   `json:"restoreSettings"`
	LastTheme       string `json:"lastTheme,omitempty"`
}

type Prefs struct {
	Dirs Dirs     `json:"dirs"`
	App  AppPrefs `json:"app"`
}

// AppDir returns the root ansizalizer config directory.
func AppDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "ansizalizer")
}

// SettingsDir returns the path to the settings subdirectory.
func SettingsDir() string {
	return filepath.Join(AppDir(), "settings")
}

// PalettesDir returns the path to the palettes subdirectory.
func PalettesDir() string {
	return filepath.Join(AppDir(), "palettes")
}

func prefsPath() (string, error) {
	dir := AppDir()
	if dir == "" {
		return "", os.ErrNotExist
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, prefsFileName), nil
}

func Load() Prefs {
	path, err := prefsPath()
	if err != nil {
		return Prefs{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Prefs{}
	}
	var p Prefs
	if err := json.Unmarshal(data, &p); err != nil {
		return Prefs{}
	}
	return p
}

func (p Prefs) Save() {
	path, err := prefsPath()
	if err != nil {
		return
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}
