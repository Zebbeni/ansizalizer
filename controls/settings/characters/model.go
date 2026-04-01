package characters

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/debug"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Ascii State = iota
	Unicode
	Custom
	AsciiAz
	AsciiNums
	AsciiSpec
	AsciiAll
	UnicodeFull
	UnicodeHalf
	UnicodeQuart
	UnicodeShadeLight
	UnicodeShadeMed
	UnicodeShadeHeavy
	SymbolsForm
	DarkVariance
	LightVariance
	Sequence
	Random
	SeedForm
	ThresholdForm
)

type Model struct {
	focus         State
	active        State
	mode          State
	charControls  State
	unicodeMode   State
	asciiMode     State
	selectionMode  State
	customInput    textinput.Model
	seedInput      numberinput.Model
	thresholdInput numberinput.Model
	ShouldClose   bool
	IsActive      bool
	width         int
}

func New(w int) Model {
	return Model{
		focus:         Unicode,
		active:        Unicode,
		mode:          Unicode,
		charControls:  Unicode,
		asciiMode:     AsciiAz,
		unicodeMode:   UnicodeHalf,
		selectionMode: DarkVariance,
		customInput:    newInput("Symbols ", "/%A"),
		seedInput:      newNumberInput("Seed ", false, 42, nil),
		thresholdInput: newNumberInput("Threshold ", true, 0, numberinput.FloatPtr(0)),
		ShouldClose:   false,
		IsActive:      false,
		width:         w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case SymbolsForm:
		if m.customInput.Focused() {
			return m.handleSymbolsFormUpdate(msg)
		}
	case SeedForm:
		if m.seedInput.Focused() {
			return m.handleSeedFormUpdate(msg)
		}
	case ThresholdForm:
		if m.thresholdInput.Focused() {
			return m.handleThresholdFormUpdate(msg)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.drawCharTabs()
}

// Selected returns the mode, charMode, and the current set of custom-defined characters
func (m Model) Selected() (State, State, []rune) {
	var charMode State

	switch m.mode {
	case Unicode:
		charMode = m.unicodeMode
	case Ascii:
		charMode = m.asciiMode
	case Custom:
		charMode = Custom
	}

	return m.mode, charMode, []rune(m.customInput.Value())
}

func (m Model) SelectionMode() State {
	return m.selectionMode
}

func (m Model) AsciiMode() State {
	return m.asciiMode
}

func (m Model) UnicodeMode() State {
	return m.unicodeMode
}

func (m *Model) SetRandomSeed(seed int64) {
	m.seedInput.SetInt64Value(seed)
}

func (m Model) VarianceThreshold() float64 {
	return m.thresholdInput.FloatValue()
}

func (m *Model) SetVarianceThreshold(t float64) {
	m.thresholdInput.SetFloatValue(t, 2)
}

func (m Model) RandomSeed() int64 {
	return m.seedInput.Int64Value()
}

func (m *Model) SetConfig(mode, asciiMode, unicodeMode, selectionMode State, customChars string) {
	debug.Log("Characters.SetConfig: mode=%d asciiMode=%d unicodeMode=%d", mode, asciiMode, unicodeMode)
	m.mode = mode
	m.charControls = mode
	m.asciiMode = asciiMode
	m.unicodeMode = unicodeMode
	m.selectionMode = selectionMode
	m.customInput.SetValue(customChars)
}

func (m *Model) ResetFocus() {
	m.focus = Unicode
	m.active = Unicode
	m.charControls = m.mode
}

func (m Model) Summary() string {
	var chars string
	switch m.mode {
	case Ascii:
		chars = "Chars: " + stateNames[m.asciiMode]
	case Unicode:
		chars = "Chars: " + stateNames[m.unicodeMode]
	case Custom:
		val := m.customInput.Value()
		if len(val) > 6 {
			val = val[:6] + ".."
		}
		chars = "Chars: " + val
	}

	summary := chars
	if m.mode == Custom {
		summary += "\nMode: " + stateNames[m.selectionMode]
	}
	return summary
}
