package characters

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/io"
)

type State int

const (
	Ascii State = iota
	Unicode
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
}

func New() Model {
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
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, io.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, io.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, io.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	colorsButtons := m.drawColorsButtons()
	modeButtons := m.drawModeButtons()
	charButtons := m.drawCharButtons()
	return lipgloss.JoinVertical(lipgloss.Top, colorsButtons, modeButtons, charButtons)
}

func (m Model) Selected() (State, State, State) {
	charMode := m.asciiMode
	if m.mode == Unicode {
		charMode = m.unicodeMode
	}
	return m.mode, charMode, m.useFgBg
}
