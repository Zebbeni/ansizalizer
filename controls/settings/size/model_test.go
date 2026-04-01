package size

import (
	"strings"
	"testing"

	"github.com/Zebbeni/ansizalizer/style"
)

func TestMain(m *testing.M) {
	style.SetTheme(style.LightOnDark)
	m.Run()
}

func TestNewDefaults(t *testing.T) {
	m := New()
	mode, width, height, charRatio := m.Info()
	if mode != Fit {
		t.Errorf("mode = %v, want Fit", mode)
	}
	if width != 150 {
		t.Errorf("width = %d, want 150", width)
	}
	if height != 50 {
		t.Errorf("height = %d, want 50", height)
	}
	if charRatio != DEFAULT_CHAR_W_TO_H_RATIO {
		t.Errorf("charRatio = %f, want %f", charRatio, DEFAULT_CHAR_W_TO_H_RATIO)
	}
}

func TestSetConfig(t *testing.T) {
	m := New()
	m.SetConfig(Fill, 80, 40, 0.5)
	mode, width, height, charRatio := m.Info()
	if mode != Fill {
		t.Errorf("mode = %v, want Fill", mode)
	}
	if width != 80 {
		t.Errorf("width = %d, want 80", width)
	}
	if height != 40 {
		t.Errorf("height = %d, want 40", height)
	}
	if charRatio != 0.5 {
		t.Errorf("charRatio = %f, want 0.5", charRatio)
	}
}

func TestSummary(t *testing.T) {
	m := New()
	m.SetConfig(Fill, 80, 40, 0.5)
	s := m.Summary()
	if !strings.Contains(s, "80") || !strings.Contains(s, "40") {
		t.Errorf("Summary() = %q, should contain dimensions", s)
	}
	if !strings.Contains(s, "Fill") {
		t.Errorf("Summary() = %q, should contain mode name", s)
	}
}
