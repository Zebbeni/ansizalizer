package characters

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

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
		Variance:       Sequence,
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
		Sequence:          Variance,
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
		Variance:       SymbolsForm,
		Sequence:          SymbolsForm,
		Random:            SymbolsForm,
	},
	Down: {
		ColorBgOff:  Ascii,
		ColorBgOn:   Ascii,
		Ascii:       AsciiAz,
		Unicode:     UnicodeShadeMed,
		Custom:      SymbolsForm,
		SymbolsForm: Variance,
	},
}

var (
	asciiCharModeMap   = map[State]bool{AsciiAz: true, AsciiNums: true, AsciiSpec: true, AsciiAll: true}
	unicodeCharModeMap = map[State]bool{UnicodeFull: true, UnicodeHalf: true, UnicodeQuart: true, UnicodeShadeLight: true, UnicodeShadeMed: true, UnicodeShadeHeavy: true}
)

func (m Model) handleThresholdFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.thresholdInput.Blur()
			m.fgBgThreshold = m.FgBgThreshold()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.thresholdInput.Blur()
		}
	}

	var cmd tea.Cmd
	m.thresholdInput, cmd = m.thresholdInput.Update(msg)
	return m, cmd
}

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
	case ThresholdForm:
		m.thresholdInput.Focus()
		return m, nil
	case SymbolsForm:
		m.mode = Custom
		m.customInput.Focus()
	case Variance, Sequence, Random:
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
		// ThresholdForm goes up to the row above it in the current tab
		if m.focus == ThresholdForm {
			switch m.charControls {
			case Ascii:
				return m.setFocus(AsciiAz)
			case Custom:
				return m.setFocus(m.selectionMode)
			}
		}
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		// Focused but inactive tabs don't navigate on down
		switch m.focus {
		case Ascii, Unicode, Custom:
			if m.charControls != m.focus {
				return m, nil
			}
		}
		// Navigate down to ThresholdForm from Ascii buttons or Custom selection modes (when visible)
		if m.colorBg {
			if _, ok := asciiCharModeMap[m.focus]; ok && m.charControls == Ascii {
				return m.setFocus(ThresholdForm)
			}
			if m.selectionMode == Variance && (m.focus == Variance || m.focus == Sequence || m.focus == Random) {
				return m.setFocus(ThresholdForm)
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
