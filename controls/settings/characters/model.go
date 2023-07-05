package characters

import (
	"github.com/charmbracelet/bubbles/key"
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
	UnicodeShadeAll
	OneColor
	TwoColor
)

type Model struct {
	focus         State
	active        State
	mode          State
	charButtons   State
	unicodeMode   State
	asciiMode     State
	useFgBg       State
	ShouldClose   bool
	ShouldUnfocus bool
	IsActive      bool
	width         int
}

func New(w int) Model {
	return Model{
		focus:         Ascii,
		active:        Ascii,
		mode:          Ascii,
		charButtons:   Ascii,
		asciiMode:     AsciiAll,
		unicodeMode:   UnicodeQuart,
		useFgBg:       OneColor,
		ShouldClose:   false,
		ShouldUnfocus: false,
		IsActive:      false,
		width:         w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) Selected() (State, State, State) {
	charMode := m.asciiMode
	if m.mode == Unicode {
		charMode = m.unicodeMode
	}
	return m.mode, charMode, m.useFgBg
}
