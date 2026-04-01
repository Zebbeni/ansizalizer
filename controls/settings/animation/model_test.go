package animation

import (
	"strings"
	"testing"
	"time"

	"github.com/Zebbeni/ansizalizer/style"
)

func TestMain(m *testing.M) {
	style.SetTheme(style.LightOnDark)
	m.Run()
}

func TestNewDefaults(t *testing.T) {
	m := New()
	if got := m.DelayMs(); got != 100 {
		t.Errorf("DelayMs() = %d, want 100", got)
	}
	if got := m.Delay(); got != 100*time.Millisecond {
		t.Errorf("Delay() = %v, want 100ms", got)
	}
}

func TestSetDelayMs(t *testing.T) {
	m := New()
	m.SetDelayMs(200)
	if got := m.DelayMs(); got != 200 {
		t.Errorf("DelayMs() = %d, want 200", got)
	}
	if got := m.Delay(); got != 200*time.Millisecond {
		t.Errorf("Delay() = %v, want 200ms", got)
	}
}

func TestDelayClampsToMin(t *testing.T) {
	m := New()
	m.SetDelayMs(1) // below min of 10
	if got := m.DelayMs(); got != 10 {
		t.Errorf("DelayMs() = %d, want 10 (clamped to min)", got)
	}
}

func TestDelayClampsToMax(t *testing.T) {
	m := New()
	m.SetDelayMs(9999) // above max of 2000
	if got := m.DelayMs(); got != 2000 {
		t.Errorf("DelayMs() = %d, want 2000 (clamped to max)", got)
	}
}

func TestSummary(t *testing.T) {
	m := New()
	m.SetDelayMs(150)
	s := m.Summary()
	if !strings.Contains(s, "150") {
		t.Errorf("Summary() = %q, should contain delay value", s)
	}
	if !strings.Contains(s, "ms") {
		t.Errorf("Summary() = %q, should contain 'ms'", s)
	}
}
