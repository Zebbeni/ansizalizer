package prefs

import (
	"encoding/json"
	"testing"
)

func TestPrefsJSONRoundTrip(t *testing.T) {
	p := Prefs{
		Dirs: Dirs{
			Browse:       "/home/user/images",
			ExportSource: "/home/user/source",
			ExportDest:   "/home/user/dest",
			SaveDir:      "/home/user/settings",
			LoadFile:     "/home/user/settings",
			PaletteLoad:  "/home/user/palettes",
		},
		App: AppPrefs{
			ShowHelp:        true,
			RestoreTheme:    true,
			RestoreSettings: false,
			LastTheme:       "Dark on Light",
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var loaded Prefs
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if loaded.Dirs.Browse != p.Dirs.Browse {
		t.Errorf("Dirs.Browse: got %q, want %q", loaded.Dirs.Browse, p.Dirs.Browse)
	}
	if loaded.Dirs.ExportSource != p.Dirs.ExportSource {
		t.Errorf("Dirs.ExportSource: got %q, want %q", loaded.Dirs.ExportSource, p.Dirs.ExportSource)
	}
	if loaded.App.ShowHelp != p.App.ShowHelp {
		t.Errorf("App.ShowHelp: got %v, want %v", loaded.App.ShowHelp, p.App.ShowHelp)
	}
	if loaded.App.RestoreTheme != p.App.RestoreTheme {
		t.Errorf("App.RestoreTheme: got %v, want %v", loaded.App.RestoreTheme, p.App.RestoreTheme)
	}
	if loaded.App.RestoreSettings != p.App.RestoreSettings {
		t.Errorf("App.RestoreSettings: got %v, want %v", loaded.App.RestoreSettings, p.App.RestoreSettings)
	}
	if loaded.App.LastTheme != p.App.LastTheme {
		t.Errorf("App.LastTheme: got %q, want %q", loaded.App.LastTheme, p.App.LastTheme)
	}
}

func TestEmptyJSONProducesZeroPrefs(t *testing.T) {
	var p Prefs
	if err := json.Unmarshal([]byte("{}"), &p); err != nil {
		t.Fatalf("Unmarshal empty JSON: %v", err)
	}
	if p.Dirs.Browse != "" {
		t.Errorf("empty JSON should produce empty Dirs.Browse, got %q", p.Dirs.Browse)
	}
	if p.App.ShowHelp != false {
		t.Error("empty JSON should produce ShowHelp=false")
	}
}

func TestAppPrefsDefaults(t *testing.T) {
	var p AppPrefs
	if p.ShowHelp {
		t.Error("default ShowHelp should be false")
	}
	if p.RestoreTheme {
		t.Error("default RestoreTheme should be false")
	}
	if p.RestoreSettings {
		t.Error("default RestoreSettings should be false")
	}
	if p.LastTheme != "" {
		t.Errorf("default LastTheme should be empty, got %q", p.LastTheme)
	}
}
