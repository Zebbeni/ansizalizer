package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
)

func TestMain(m *testing.M) {
	style.SetTheme(style.LightOnDark)
	os.Exit(m.Run())
}

func TestDefaultConfigHasRequiredFields(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Size.Width == 0 {
		t.Error("DefaultConfig().Size.Width should not be 0")
	}
	if cfg.Size.Height == 0 {
		t.Error("DefaultConfig().Size.Height should not be 0")
	}
	if cfg.Size.CharRatio == 0 {
		t.Error("DefaultConfig().Size.CharRatio should not be 0")
	}
	if cfg.Characters.Mode == "" {
		t.Error("DefaultConfig().Characters.Mode should not be empty")
	}
	if cfg.Characters.UnicodeMode == "" {
		t.Error("DefaultConfig().Characters.UnicodeMode should not be empty")
	}
	if cfg.Characters.AsciiMode == "" {
		t.Error("DefaultConfig().Characters.AsciiMode should not be empty")
	}
}

func TestConfigJSONRoundTrip(t *testing.T) {
	cfg := DefaultConfig()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	// Check key fields survive round-trip
	if loaded.Size.Width != cfg.Size.Width {
		t.Errorf("Size.Width: got %d, want %d", loaded.Size.Width, cfg.Size.Width)
	}
	if loaded.Size.Height != cfg.Size.Height {
		t.Errorf("Size.Height: got %d, want %d", loaded.Size.Height, cfg.Size.Height)
	}
	if loaded.Characters.Mode != cfg.Characters.Mode {
		t.Errorf("Characters.Mode: got %q, want %q", loaded.Characters.Mode, cfg.Characters.Mode)
	}
	if loaded.Colors.UseTrueColor != cfg.Colors.UseTrueColor {
		t.Errorf("Colors.UseTrueColor: got %v, want %v", loaded.Colors.UseTrueColor, cfg.Colors.UseTrueColor)
	}
	if loaded.Alpha.AlphaThreshold != cfg.Alpha.AlphaThreshold {
		t.Errorf("Alpha.AlphaThreshold: got %f, want %f", loaded.Alpha.AlphaThreshold, cfg.Alpha.AlphaThreshold)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	cfg := DefaultConfig()
	cfg.Size.Width = 123
	cfg.Size.Height = 456
	cfg.Characters.Mode = "Ascii"

	if err := SaveConfig(cfg, path); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if loaded.Size.Width != 123 {
		t.Errorf("Size.Width: got %d, want 123", loaded.Size.Width)
	}
	if loaded.Size.Height != 456 {
		t.Errorf("Size.Height: got %d, want 456", loaded.Size.Height)
	}
	if loaded.Characters.Mode != "Ascii" {
		t.Errorf("Characters.Mode: got %q, want %q", loaded.Characters.Mode, "Ascii")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/settings.json")
	if err == nil {
		t.Error("LoadConfig should return error for missing file")
	}
}

func TestCharModeNameMapsAreBijective(t *testing.T) {
	for k, v := range charModeNames {
		if rev, ok := charModeFromName[v]; !ok {
			t.Errorf("charModeFromName missing key %q (from charModeNames[%d])", v, k)
		} else if rev != k {
			t.Errorf("charModeFromName[%q] = %d, want %d", v, rev, k)
		}
	}
}

func TestAsciiModeNameMapsAreBijective(t *testing.T) {
	for k, v := range asciiModeNames {
		if rev, ok := asciiModeFromName[v]; !ok {
			t.Errorf("asciiModeFromName missing key %q", v)
		} else if rev != k {
			t.Errorf("asciiModeFromName[%q] = %d, want %d", v, rev, k)
		}
	}
}

func TestUnicodeModeNameMapsAreBijective(t *testing.T) {
	for k, v := range unicodeModeNames {
		if rev, ok := unicodeModeFromName[v]; !ok {
			t.Errorf("unicodeModeFromName missing key %q", v)
		} else if rev != k {
			t.Errorf("unicodeModeFromName[%q] = %d, want %d", v, rev, k)
		}
	}
}

func TestApplyConfigRoundTrip(t *testing.T) {
	m := New(30, prefs.Dirs{})
	cfg := DefaultConfig()
	m.ApplyConfig(cfg)
	exported := m.ExportConfig()

	if exported.Size.Width != cfg.Size.Width {
		t.Errorf("Size.Width: got %d, want %d", exported.Size.Width, cfg.Size.Width)
	}
	if exported.Characters.Mode != cfg.Characters.Mode {
		t.Errorf("Characters.Mode: got %q, want %q", exported.Characters.Mode, cfg.Characters.Mode)
	}
	if exported.Characters.UnicodeMode != cfg.Characters.UnicodeMode {
		t.Errorf("Characters.UnicodeMode: got %q, want %q", exported.Characters.UnicodeMode, cfg.Characters.UnicodeMode)
	}
	if exported.Alpha.UseAlpha != cfg.Alpha.UseAlpha {
		t.Errorf("Alpha.UseAlpha: got %v, want %v", exported.Alpha.UseAlpha, cfg.Alpha.UseAlpha)
	}
}
