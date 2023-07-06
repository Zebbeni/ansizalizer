package characters

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {
		Ascii:             Unicode,
		Unicode:           Custom,
		AsciiAz:           AsciiNums,
		AsciiNums:         AsciiSpec,
		AsciiSpec:         AsciiAll,
		UnicodeFull:       UnicodeHalf,
		UnicodeHalf:       UnicodeQuart,
		UnicodeQuart:      UnicodeShadeLight,
		UnicodeShadeLight: UnicodeShadeMed,
		UnicodeShadeMed:   UnicodeShadeHeavy,
		TwoColor:          OneColor,
	},
	Left: {
		Unicode:           Ascii,
		Custom:            Unicode,
		AsciiAll:          AsciiSpec,
		AsciiSpec:         AsciiNums,
		AsciiNums:         AsciiAz,
		UnicodeShadeHeavy: UnicodeShadeMed,
		UnicodeShadeMed:   UnicodeShadeLight,
		UnicodeShadeLight: UnicodeQuart,
		UnicodeQuart:      UnicodeHalf,
		UnicodeHalf:       UnicodeFull,
		OneColor:          TwoColor,
	},
	Up: {
		Ascii:             TwoColor,
		Unicode:           OneColor,
		Custom:            OneColor,
		AsciiAz:           Ascii,
		AsciiNums:         Ascii,
		AsciiSpec:         Ascii,
		AsciiAll:          Ascii,
		UnicodeFull:       Unicode,
		UnicodeHalf:       Unicode,
		UnicodeQuart:      Unicode,
		UnicodeShadeLight: Unicode,
		UnicodeShadeMed:   Unicode,
		UnicodeShadeHeavy: Unicode,
		SymbolsForm:       Custom,
	},
	Down: {
		OneColor: Unicode,
		TwoColor: Ascii,
		Ascii:    AsciiAz,
		Unicode:  UnicodeShadeMed,
		Custom:   SymbolsForm,
	},
}

var (
	asciiCharModeMap   = map[State]bool{AsciiAz: true, AsciiNums: true, AsciiSpec: true, AsciiAll: true}
	unicodeCharModeMap = map[State]bool{UnicodeFull: true, UnicodeHalf: true, UnicodeQuart: true, UnicodeShadeLight: true, UnicodeShadeMed: true, UnicodeShadeHeavy: true}
)

func (m Model) handleSymbolsFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.customInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.customInput.Blur()
		}
	}

	var cmd tea.Cmd
	m.customInput, cmd = m.customInput.Update(msg)
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus

	switch m.active {
	case Ascii:
		m.mode = Ascii
	case Unicode:
		m.mode = Unicode
	case Custom:
		m.mode = Custom
	case SymbolsForm:
		m.customInput.Focus()
	case OneColor, TwoColor:
		m.useFgBg = m.active
	default:
		switch m.charButtons {
		case Ascii:
			if _, ok := asciiCharModeMap[m.active]; ok {
				m.asciiMode = m.active
				m.mode = Ascii
			}
		case Unicode:
			if _, ok := unicodeCharModeMap[m.active]; ok {
				m.unicodeMode = m.active
				m.mode = Unicode
			}
		}
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	}
	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	switch m.focus {
	case Ascii:
		m.charButtons = Ascii
		m.mode = Ascii
	case Unicode:
		m.charButtons = Unicode
		m.mode = Unicode
	case Custom:
		m.charButtons = Custom
		m.mode = Custom
	}
	return m, nil
}
