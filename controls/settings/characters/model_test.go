package characters

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
	mode, charMode, _ := m.Selected()
	if mode != Unicode {
		t.Errorf("mode = %v, want Unicode", mode)
	}
	if charMode != UnicodeHalf {
		t.Errorf("charMode = %v, want UnicodeHalf", charMode)
	}
	if m.AsciiMode() != AsciiAz {
		t.Errorf("AsciiMode() = %v, want AsciiAz", m.AsciiMode())
	}
	if m.UnicodeMode() != UnicodeHalf {
		t.Errorf("UnicodeMode() = %v, want UnicodeHalf", m.UnicodeMode())
	}
	if m.SelectionMode() != DarkVariance {
		t.Errorf("SelectionMode() = %v, want DarkVariance", m.SelectionMode())
	}
}

func TestNewDefaultSeed(t *testing.T) {
	m := New(30)
	if got := m.RandomSeed(); got != 42 {
		t.Errorf("RandomSeed() = %d, want 42", got)
	}
}

func TestNewDefaultThreshold(t *testing.T) {
	m := New(30)
	if got := m.VarianceThreshold(); got != 0 {
		t.Errorf("VarianceThreshold() = %f, want 0", got)
	}
}

func TestSetConfig(t *testing.T) {
	m := New(30)
	m.SetConfig(Ascii, AsciiNums, UnicodeQuart, LightVariance, "ABC")
	mode, charMode, customChars := m.Selected()
	if mode != Ascii {
		t.Errorf("mode = %v, want Ascii", mode)
	}
	if charMode != AsciiNums {
		t.Errorf("charMode = %v, want AsciiNums", charMode)
	}
	if m.UnicodeMode() != UnicodeQuart {
		t.Errorf("UnicodeMode() = %v, want UnicodeQuart", m.UnicodeMode())
	}
	if m.SelectionMode() != LightVariance {
		t.Errorf("SelectionMode() = %v, want LightVariance", m.SelectionMode())
	}
	if string(customChars) != "ABC" {
		t.Errorf("customChars = %q, want %q", string(customChars), "ABC")
	}
}

func TestSetConfigUnicodeMode(t *testing.T) {
	m := New(30)
	m.SetConfig(Unicode, AsciiAz, UnicodeFull, DarkVariance, "")
	mode, charMode, _ := m.Selected()
	if mode != Unicode {
		t.Errorf("mode = %v, want Unicode", mode)
	}
	if charMode != UnicodeFull {
		t.Errorf("charMode = %v, want UnicodeFull", charMode)
	}
}

func TestSetRandomSeed(t *testing.T) {
	m := New(30)
	m.SetRandomSeed(12345)
	if got := m.RandomSeed(); got != 12345 {
		t.Errorf("RandomSeed() = %d, want 12345", got)
	}
}

func TestSetVarianceThreshold(t *testing.T) {
	m := New(30)
	m.SetVarianceThreshold(0.75)
	if got := m.VarianceThreshold(); got != 0.75 {
		t.Errorf("VarianceThreshold() = %f, want 0.75", got)
	}
}

func TestSummaryUnicode(t *testing.T) {
	m := New(30)
	s := m.Summary()
	if !strings.Contains(s, "Chars:") {
		t.Errorf("Summary() = %q, should contain 'Chars:'", s)
	}
}

func TestSummaryAscii(t *testing.T) {
	m := New(30)
	m.SetConfig(Ascii, AsciiAll, UnicodeHalf, DarkVariance, "")
	s := m.Summary()
	if !strings.Contains(s, "All") {
		t.Errorf("Summary() = %q, should contain 'All' for AsciiAll mode", s)
	}
}

func TestSummaryCustom(t *testing.T) {
	m := New(30)
	m.SetConfig(Custom, AsciiAz, UnicodeHalf, DarkVariance, "ABCDEF")
	s := m.Summary()
	if !strings.Contains(s, "ABCDEF") {
		t.Errorf("Summary() = %q, should contain custom chars", s)
	}
}

func TestSummaryCustomTruncatesLong(t *testing.T) {
	m := New(30)
	m.SetConfig(Custom, AsciiAz, UnicodeHalf, DarkVariance, "ABCDEFGHIJ")
	s := m.Summary()
	if !strings.Contains(s, "..") {
		t.Errorf("Summary() = %q, should truncate long custom chars with '..'", s)
	}
}
