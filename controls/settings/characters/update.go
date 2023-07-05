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
	Right: {Ascii: Unicode,
		AsciiAz:           AsciiNums,
		AsciiNums:         AsciiSpec,
		AsciiSpec:         AsciiAll,
		UnicodeFull:       UnicodeHalf,
		UnicodeHalf:       UnicodeQuart,
		UnicodeQuart:      UnicodeShadeLight,
		UnicodeShadeLight: UnicodeShadeMed,
		UnicodeShadeMed:   UnicodeShadeHeavy,
		UnicodeShadeHeavy: UnicodeShadeAll,
		TwoColor:          OneColor,
	},
	Left: {Unicode: Ascii,
		AsciiAll:          AsciiSpec,
		AsciiSpec:         AsciiNums,
		AsciiNums:         AsciiAz,
		UnicodeShadeAll:   UnicodeShadeHeavy,
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
		UnicodeShadeAll:   Unicode,
	},
	Down: {
		OneColor: Unicode,
		TwoColor: Ascii,

		Ascii:   AsciiAz,
		Unicode: UnicodeShadeMed,
	},
}

var (
	asciiCharModeMap   = map[State]bool{AsciiAz: true, AsciiNums: true, AsciiSpec: true, AsciiAll: true}
	unicodeCharModeMap = map[State]bool{UnicodeFull: true, UnicodeHalf: true, UnicodeQuart: true, UnicodeShadeLight: true, UnicodeShadeMed: true, UnicodeShadeHeavy: true, UnicodeShadeAll: true}
)

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
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		}
	}
	return m, cmd
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	if m.focus == Ascii {
		m.charButtons = Ascii
	} else if m.focus == Unicode {
		m.charButtons = Unicode
	}
	return m, nil
}
