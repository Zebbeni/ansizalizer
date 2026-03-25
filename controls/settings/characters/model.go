package characters

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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
	DarkToLight
	Sequence
	Random
	ColorBgOff
	ColorBgOn
)

type Model struct {
	focus         State
	active        State
	mode          State
	charControls  State
	unicodeMode   State
	asciiMode     State
	colorBg       bool
	selectionMode State
	customInput   textinput.Model
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
		colorBg:       true,
		selectionMode: DarkToLight,
		customInput:   newInput("Symbols", "/%A"),
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
	colorsButtons := m.drawColorsButtons()
	charTabs := m.drawCharTabs()
	return lipgloss.JoinVertical(lipgloss.Top, colorsButtons, charTabs)
}

// Selected returns the mode, charMode, whether to use background color, and the
// current set of custom-defined characters
func (m Model) Selected() (State, State, bool, []rune) {
	var charMode State

	switch m.mode {
	case Unicode:
		charMode = m.unicodeMode
	case Ascii:
		charMode = m.asciiMode
	case Custom:
		charMode = Custom
	}

	return m.mode, charMode, m.colorBg, []rune(m.customInput.Value())
}

func (m Model) SelectionMode() State {
	return m.selectionMode
}

func (m *Model) SetConfig(mode, asciiMode, unicodeMode, selectionMode State, colorBg bool, customChars string) {
	m.mode = mode
	m.charControls = mode
	m.asciiMode = asciiMode
	m.unicodeMode = unicodeMode
	m.colorBg = colorBg
	m.selectionMode = selectionMode
	m.customInput.SetValue(customChars)
}

func (m *Model) ResetFocus() {
	m.focus = Unicode
	m.active = Unicode
	m.charControls = m.mode
}

func (m Model) Summary() string {
	colors := "Colors: FG Only"
	if m.colorBg {
		colors = "Colors: FG + BG"
	}

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

	summary := colors + "\n" + chars
	if m.mode == Custom {
		summary += "\nMode: " + stateNames[m.selectionMode]
	}
	return summary
}
