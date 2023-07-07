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
	OneColor
	TwoColor
)

type Model struct {
	focus        State
	active       State
	mode         State
	charControls State
	unicodeMode  State
	asciiMode    State
	useFgBg      State
	customInput  textinput.Model
	ShouldClose  bool
	IsActive     bool
	width        int
}

func New(w int) Model {
	return Model{
		focus:        Ascii,
		active:       Ascii,
		mode:         Ascii,
		charControls: Ascii,
		asciiMode:    AsciiAll,
		unicodeMode:  UnicodeFull,
		useFgBg:      OneColor,
		customInput:  newInput("Symbols", "/%A"),
		ShouldClose:  false,
		IsActive:     false,
		width:        w,
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

// Selected returns the mode, charMode, whether to use two colors, and the
// current set of custom-defined characters
func (m Model) Selected() (State, State, State, []rune) {
	var charMode State

	switch m.mode {
	case Unicode:
		charMode = m.unicodeMode
	case Ascii:
		charMode = m.asciiMode
	case Custom:
		charMode = Custom
	}

	return m.mode, charMode, m.useFgBg, []rune(m.customInput.Value())
}
