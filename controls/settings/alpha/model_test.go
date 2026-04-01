package alpha

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
	if !m.ShouldOutputAlpha() {
		t.Error("ShouldOutputAlpha() should be true by default")
	}
	if m.TrimAlpha() {
		t.Error("TrimAlpha() should be false by default")
	}
	if got := m.AlphaThreshold(); got != 0.5 {
		t.Errorf("AlphaThreshold() = %f, want 0.5", got)
	}
}

func TestSetConfig(t *testing.T) {
	m := New(30)
	m.SetConfig(false, true)
	if m.ShouldOutputAlpha() {
		t.Error("ShouldOutputAlpha() should be false after SetConfig(false, true)")
	}
	if !m.TrimAlpha() {
		t.Error("TrimAlpha() should be true after SetConfig(false, true)")
	}
}

func TestSetAlphaThreshold(t *testing.T) {
	m := New(30)
	m.SetAlphaThreshold(0.75)
	if got := m.AlphaThreshold(); got != 0.75 {
		t.Errorf("AlphaThreshold() = %f, want 0.75", got)
	}
}

func TestAlphaThresholdClamped(t *testing.T) {
	m := New(30)
	m.SetAlphaThreshold(1.5)
	// The numberinput clamps to max=1.0 on blur, but SetFloatValue doesn't auto-clamp
	// FloatValue() does clamp on read
	if got := m.AlphaThreshold(); got > 1.0 {
		t.Errorf("AlphaThreshold() = %f, should be clamped to 1.0", got)
	}
}

func TestSummaryAlphaOff(t *testing.T) {
	m := New(30)
	m.SetConfig(false, false)
	s := m.Summary()
	if strings.Contains(s, "Transparent") {
		t.Errorf("Summary() = %q, should not contain 'Transparent' when alpha is off", s)
	}
}

func TestSummaryAlphaOn(t *testing.T) {
	m := New(30)
	m.SetConfig(true, false)
	s := m.Summary()
	if !strings.Contains(s, "Transparent") {
		t.Errorf("Summary() = %q, should contain 'Transparent' when alpha is on", s)
	}
}

func TestSummaryTrim(t *testing.T) {
	m := New(30)
	m.SetConfig(true, true)
	s := m.Summary()
	if !strings.Contains(s, "Trim") {
		t.Errorf("Summary() = %q, should contain 'Trim'", s)
	}
}

func TestSummaryThreshold(t *testing.T) {
	m := New(30)
	s := m.Summary()
	if !strings.Contains(s, "Thresh") {
		t.Errorf("Summary() = %q, should contain 'Thresh'", s)
	}
}
