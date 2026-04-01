package textstyle

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
	m := New(30)
	if m.Bold() {
		t.Error("Bold() should be false by default")
	}
	if m.Italic() {
		t.Error("Italic() should be false by default")
	}
	if m.Underline() {
		t.Error("Underline() should be false by default")
	}
	if m.Strikethrough() {
		t.Error("Strikethrough() should be false by default")
	}
}

func TestSetConfig(t *testing.T) {
	m := New(30)
	m.SetConfig(true, false, true, false)
	if !m.Bold() {
		t.Error("Bold() should be true after SetConfig")
	}
	if m.Italic() {
		t.Error("Italic() should be false after SetConfig")
	}
	if !m.Underline() {
		t.Error("Underline() should be true after SetConfig")
	}
	if m.Strikethrough() {
		t.Error("Strikethrough() should be false after SetConfig")
	}
}

func TestSummaryAllOff(t *testing.T) {
	m := New(30)
	s := m.Summary()
	if !strings.Contains(s, "None") {
		t.Errorf("Summary() = %q, want 'None' when all styles are off", s)
	}
}

func TestSummaryWithStyles(t *testing.T) {
	m := New(30)
	m.SetConfig(true, true, false, false)
	s := m.Summary()
	if !strings.Contains(s, "Bold") {
		t.Errorf("Summary() = %q, should contain 'Bold'", s)
	}
	if !strings.Contains(s, "Italic") {
		t.Errorf("Summary() = %q, should contain 'Italic'", s)
	}
}

func TestSummaryAllOn(t *testing.T) {
	m := New(30)
	m.SetConfig(true, true, true, true)
	s := m.Summary()
	for _, style := range []string{"Bold", "Italic", "Underline", "Strikethrough"} {
		if !strings.Contains(s, style) {
			t.Errorf("Summary() = %q, should contain %q", s, style)
		}
	}
}
