package adjust

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
	if got := m.Brightness(); got != 0 {
		t.Errorf("Brightness() = %d, want 0", got)
	}
	if got := m.Contrast(); got != 0 {
		t.Errorf("Contrast() = %d, want 0", got)
	}
}

func TestSetConfig(t *testing.T) {
	m := New(30)
	m.SetConfig(50, -30)
	if got := m.Brightness(); got != 50 {
		t.Errorf("Brightness() = %d, want 50", got)
	}
	if got := m.Contrast(); got != -30 {
		t.Errorf("Contrast() = %d, want -30", got)
	}
}

func TestSummary(t *testing.T) {
	m := New(30)
	m.SetConfig(25, -10)
	s := m.Summary()
	if !strings.Contains(s, "25") {
		t.Errorf("Summary() = %q, should contain brightness value", s)
	}
	if !strings.Contains(s, "-10") {
		t.Errorf("Summary() = %q, should contain contrast value", s)
	}
}
