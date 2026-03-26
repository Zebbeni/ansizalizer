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
		ColorBgOff:          ColorBgOn,
		DarkToLight:       Sequence,
		Sequence:          Random,
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
		ColorBgOn:          ColorBgOff,
		Random:            Sequence,
		Sequence:          DarkToLight,
	},
	Up: {
		Ascii:             ColorBgOff,
		Unicode:           ColorBgOff,
		Custom:            ColorBgOff,
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
		DarkToLight:       SymbolsForm,
		Sequence:          SymbolsForm,
		Random:            SymbolsForm,
	},
	Down: {
		ColorBgOff:    Custom,
		ColorBgOn:    Custom,
		Ascii:       AsciiAz,
		Unicode:     UnicodeShadeMed,
		Custom:      SymbolsForm,
		SymbolsForm: DarkToLight,
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
		m.charControls = Ascii
	case Unicode:
		m.mode = Unicode
		m.charControls = Unicode
	case Custom:
		m.mode = Custom
		m.charControls = Custom
	case SymbolsForm:
		m.mode = Custom
		m.customInput.Focus()
	case DarkToLight, Sequence, Random:
		m.selectionMode = m.active
	case ColorBgOff:
		m.colorBg = false
	case ColorBgOn:
		m.colorBg = true
	default:
		switch m.charControls {
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
		// Activate tab if focused but not active before navigating into its content
		switch m.focus {
		case Ascii, Unicode, Custom:
			if m.charControls != m.focus {
				m.charControls = m.focus
				m.mode = m.focus
			}
		}
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
	return m, nil
}
