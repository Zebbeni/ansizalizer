package numberinput

import (
	"testing"
)

func TestNewInt(t *testing.T) {
	m := New(Options{
		Prompt:  "Width ",
		IsFloat: false,
		Default: 50,
	})
	if got := m.IntValue(); got != 50 {
		t.Errorf("IntValue() = %d, want 50", got)
	}
	if got := m.Value(); got != "50" {
		t.Errorf("Value() = %q, want %q", got, "50")
	}
}

func TestNewFloat(t *testing.T) {
	m := New(Options{
		Prompt:  "Threshold ",
		IsFloat: true,
		Default: 0.5,
	})
	if got := m.FloatValue(); got != 0.5 {
		t.Errorf("FloatValue() = %f, want 0.5", got)
	}
}

func TestClampMin(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Min:     FloatPtr(-100),
		Max:     FloatPtr(100),
		Default: 0,
	})
	m.SetValue("-200")
	if got := m.IntValue(); got != -100 {
		t.Errorf("IntValue() = %d, want -100 (clamped to min)", got)
	}
}

func TestClampMax(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Min:     FloatPtr(-100),
		Max:     FloatPtr(100),
		Default: 0,
	})
	m.SetValue("999")
	if got := m.IntValue(); got != 100 {
		t.Errorf("IntValue() = %d, want 100 (clamped to max)", got)
	}
}

func TestClampFloat(t *testing.T) {
	m := New(Options{
		IsFloat: true,
		Min:     FloatPtr(0),
		Max:     FloatPtr(1),
		Default: 0.5,
	})
	m.SetValue("1.5")
	if got := m.FloatValue(); got != 1.0 {
		t.Errorf("FloatValue() = %f, want 1.0 (clamped to max)", got)
	}
	m.SetValue("-0.5")
	if got := m.FloatValue(); got != 0.0 {
		t.Errorf("FloatValue() = %f, want 0.0 (clamped to min)", got)
	}
}

func TestFallbackOnParseError(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Default: 42,
	})
	m.SetValue("abc")
	if got := m.IntValue(); got != 42 {
		t.Errorf("IntValue() = %d, want 42 (fallback)", got)
	}
}

func TestFallbackFloat(t *testing.T) {
	m := New(Options{
		IsFloat: true,
		Default: 1.0,
	})
	m.SetValue("not-a-number")
	if got := m.FloatValue(); got != 1.0 {
		t.Errorf("FloatValue() = %f, want 1.0 (fallback)", got)
	}
}

func TestFloat32Value(t *testing.T) {
	m := New(Options{
		IsFloat: true,
		Default: 0.75,
	})
	if got := m.Float32Value(); got != 0.75 {
		t.Errorf("Float32Value() = %f, want 0.75", got)
	}
}

func TestInt64Value(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Default: 42,
	})
	m.SetInt64Value(9999999999)
	if got := m.Int64Value(); got != 9999999999 {
		t.Errorf("Int64Value() = %d, want 9999999999", got)
	}
}

func TestSetIntValue(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	m.SetIntValue(77)
	if got := m.IntValue(); got != 77 {
		t.Errorf("IntValue() = %d, want 77", got)
	}
}

func TestSetFloatValue(t *testing.T) {
	m := New(Options{IsFloat: true, Default: 0})
	m.SetFloatValue(3.14, 2)
	if got := m.Value(); got != "3.14" {
		t.Errorf("Value() = %q, want %q", got, "3.14")
	}
}

func TestValidateRejectsNonNumeric(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	if err := m.validate("abc"); err == nil {
		t.Error("validate(\"abc\") should return error for int mode")
	}
}

func TestValidateAllowsEmpty(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	if err := m.validate(""); err != nil {
		t.Errorf("validate(\"\") = %v, want nil", err)
	}
}

func TestValidateAllowsMinus(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0, Min: FloatPtr(-100)})
	if err := m.validate("-"); err != nil {
		t.Errorf("validate(\"-\") = %v, want nil", err)
	}
}

func TestValidateAllowsDotForFloat(t *testing.T) {
	m := New(Options{IsFloat: true, Default: 0})
	if err := m.validate("."); err != nil {
		t.Errorf("validate(\".\") = %v, want nil for float mode", err)
	}
}

func TestValidateRejectsDecimalForInt(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	if err := m.validate("1.5"); err == nil {
		t.Error("validate(\"1.5\") should return error for int mode")
	}
}

func TestValidateRejectsAboveMax(t *testing.T) {
	m := New(Options{IsFloat: true, Max: FloatPtr(1.0), Default: 0})
	if err := m.validate("1.5"); err == nil {
		t.Error("validate(\"1.5\") should return error when max=1.0")
	}
}

func TestBlurResolvesValue(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Min:     FloatPtr(0),
		Max:     FloatPtr(100),
		Default: 50,
	})
	m.SetValue("999")
	m.Blur()
	if got := m.Value(); got != "100" {
		t.Errorf("after Blur, Value() = %q, want %q (resolved to clamped)", got, "100")
	}
}

func TestBlurResolvesInvalidToDefault(t *testing.T) {
	m := New(Options{
		IsFloat: false,
		Default: 42,
	})
	m.SetValue("xyz")
	m.Blur()
	if got := m.Value(); got != "42" {
		t.Errorf("after Blur, Value() = %q, want %q (resolved to default)", got, "42")
	}
}

func TestFloatPtr(t *testing.T) {
	p := FloatPtr(3.14)
	if p == nil || *p != 3.14 {
		t.Errorf("FloatPtr(3.14) = %v, want non-nil *3.14", p)
	}
}

func TestNoMin(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	m.SetValue("-9999")
	if got := m.IntValue(); got != -9999 {
		t.Errorf("IntValue() = %d, want -9999 (no min bound)", got)
	}
}

func TestNoMax(t *testing.T) {
	m := New(Options{IsFloat: false, Default: 0})
	m.SetValue("9999")
	if got := m.IntValue(); got != 9999 {
		t.Errorf("IntValue() = %d, want 9999 (no max bound)", got)
	}
}
