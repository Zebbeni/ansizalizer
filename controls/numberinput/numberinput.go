package numberinput

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Model wraps a bubbles textinput, restricting values to numeric input
// and optionally clamping to a min/max range.
type Model struct {
	input    textinput.Model
	isFloat  bool
	min      *float64
	max      *float64
	fallback float64
}

// Options configures a new NumberInput.
type Options struct {
	Prompt    string
	CharLimit int
	IsFloat   bool    // true for float values, false for int
	Min       *float64 // nil = no lower bound
	Max       *float64 // nil = no upper bound
	Default   float64  // initial value and fallback on parse error
}

// New creates a NumberInput from the given options.
func New(opts Options) Model {
	input := textinput.New()
	input.Prompt = opts.Prompt
	if opts.CharLimit > 0 {
		input.CharLimit = opts.CharLimit
	}

	m := Model{
		input:    input,
		isFloat:  opts.IsFloat,
		min:      opts.Min,
		max:      opts.Max,
		fallback: opts.Default,
	}

	input.Validate = m.validate
	m.input = input

	if opts.IsFloat {
		m.input.SetValue(strconv.FormatFloat(opts.Default, 'f', -1, 64))
	} else {
		m.input.SetValue(strconv.Itoa(int(opts.Default)))
	}

	return m
}

// validate checks that the current text is valid numeric input within bounds.
func (m Model) validate(s string) error {
	if s == "" || s == "-" || (m.isFloat && (s == "." || s == "-.")) {
		return nil
	}

	if m.isFloat {
		if strings.Count(s, ".") > 1 {
			return fmt.Errorf("invalid number")
		}
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("not a number")
		}
		if m.max != nil && val > *m.max {
			return fmt.Errorf("must be at most %v", *m.max)
		}
		// Allow typing negative values even when min >= 0:
		// only reject if the string is a complete number (not just "-")
		if m.min != nil && val < *m.min && !strings.HasSuffix(s, "-") {
			// Don't reject mid-typing — the user might be entering a
			// negative number or deleting digits. Clamping happens on read.
		}
	} else {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("not a number")
		}
		if strings.Contains(s, ".") {
			return fmt.Errorf("integer required")
		}
		if m.max != nil && val > *m.max {
			return fmt.Errorf("must be at most %v", *m.max)
		}
	}
	return nil
}

// IntValue returns the current value as an int, clamped to bounds.
// On parse error, returns the default.
func (m Model) IntValue() int {
	val, err := strconv.Atoi(m.input.Value())
	if err != nil {
		return int(m.fallback)
	}
	return int(m.clamp(float64(val)))
}

// FloatValue returns the current value as a float64, clamped to bounds.
// On parse error, returns the default.
func (m Model) FloatValue() float64 {
	val, err := strconv.ParseFloat(m.input.Value(), 64)
	if err != nil {
		return m.fallback
	}
	return m.clamp(val)
}

// Float32Value returns the current value as a float32, clamped to bounds.
// On parse error, returns the default.
func (m Model) Float32Value() float32 {
	return float32(m.FloatValue())
}

// Int64Value returns the current value as an int64, clamped to bounds.
// On parse error, returns the default.
func (m Model) Int64Value() int64 {
	val, err := strconv.ParseInt(m.input.Value(), 10, 64)
	if err != nil {
		return int64(m.fallback)
	}
	return int64(m.clamp(float64(val)))
}

func (m Model) clamp(val float64) float64 {
	if m.min != nil && val < *m.min {
		return *m.min
	}
	if m.max != nil && val > *m.max {
		return *m.max
	}
	return val
}

// SetValue sets the displayed value.
func (m *Model) SetValue(s string) {
	m.input.SetValue(s)
}

// SetIntValue sets the displayed value from an int.
func (m *Model) SetIntValue(v int) {
	m.input.SetValue(strconv.Itoa(v))
}

// SetFloatValue sets the displayed value from a float64.
func (m *Model) SetFloatValue(v float64, prec int) {
	m.input.SetValue(strconv.FormatFloat(v, 'f', prec, 64))
}

// SetInt64Value sets the displayed value from an int64.
func (m *Model) SetInt64Value(v int64) {
	m.input.SetValue(strconv.FormatInt(v, 10))
}

// Value returns the raw string value.
func (m Model) Value() string {
	return m.input.Value()
}

// Focus sets the input to focused state.
func (m *Model) Focus() tea.Cmd {
	return m.input.Focus()
}

// Blur removes focus and resolves the displayed value to its clamped result,
// so the user can see whether their input was rejected or bounded.
func (m *Model) Blur() {
	m.resolve()
	m.input.Blur()
}

// resolve replaces the displayed text with the parsed, clamped value.
// Invalid input is replaced with the default; values outside bounds are clamped.
func (m *Model) resolve() {
	if m.isFloat {
		v := m.FloatValue()
		m.input.SetValue(strconv.FormatFloat(v, 'f', -1, 64))
	} else {
		v := m.IntValue()
		m.input.SetValue(strconv.Itoa(v))
	}
}

// Focused returns whether the input is focused.
func (m Model) Focused() bool {
	return m.input.Focused()
}

// Update delegates to the underlying textinput.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the input.
func (m Model) View() string {
	return m.input.View()
}

// Styles returns the underlying textinput styles for customization.
func (m Model) Styles() textinput.Styles {
	return m.input.Styles()
}

// SetStyles sets the underlying textinput styles.
func (m *Model) SetStyles(s textinput.Styles) {
	m.input.SetStyles(s)
}

// SetVirtualCursor sets whether the virtual cursor is shown.
func (m *Model) SetVirtualCursor(show bool) {
	m.input.SetVirtualCursor(show)
}

// CursorEnd moves the cursor to the end of the input.
func (m *Model) CursorEnd() {
	m.input.CursorEnd()
}

// SetWidth sets the width of the input field.
func (m *Model) SetWidth(w int) {
	m.input.SetWidth(w)
}

// Prompt returns the prompt string.
func (m Model) Prompt() string {
	return m.input.Prompt
}

// Ptr helpers for setting min/max inline.
func FloatPtr(v float64) *float64 { return &v }
