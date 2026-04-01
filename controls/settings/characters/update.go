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
		Unicode:           Ascii,
		Ascii:             Custom,
		AsciiAz:           AsciiNums,
		AsciiNums:         AsciiSpec,
		AsciiSpec:         AsciiAll,
		UnicodeFull:       UnicodeHalf,
		UnicodeHalf:       UnicodeQuart,
		UnicodeQuart:      UnicodeShadeLight,
		UnicodeShadeLight: UnicodeShadeMed,
		UnicodeShadeMed:   UnicodeShadeHeavy,
		DarkVariance:      LightVariance,
		Sequence:           Random,
	},
	Left: {
		Ascii:             Unicode,
		Custom:            Ascii,
		AsciiAll:          AsciiSpec,
		AsciiSpec:         AsciiNums,
		AsciiNums:         AsciiAz,
		UnicodeShadeHeavy: UnicodeShadeMed,
		UnicodeShadeMed:   UnicodeShadeLight,
		UnicodeShadeLight: UnicodeQuart,
		UnicodeQuart:      UnicodeHalf,
		UnicodeHalf:       UnicodeFull,
		LightVariance:     DarkVariance,
		Random:            Sequence,
	},
	Up: {
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
		// DarkVariance/LightVariance Up is context-dependent — handled in handleNav
		Sequence:           DarkVariance,
		Random:             DarkVariance,
		SeedForm:           Random,
		ThresholdForm:      DarkVariance,
	},
	Down: {
		Ascii:         AsciiAz,
		Unicode:       UnicodeShadeMed,
		Custom:        SymbolsForm,
		SymbolsForm:   DarkVariance,
		AsciiAz:       DarkVariance,
		AsciiNums:     DarkVariance,
		AsciiSpec:     DarkVariance,
		AsciiAll:      DarkVariance,
		DarkVariance:  Sequence,
		LightVariance: Sequence,
		ThresholdForm: Sequence,
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

func (m Model) handleThresholdFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.thresholdInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.thresholdInput.Blur()
		}
	}

	var cmd tea.Cmd
	m.thresholdInput, cmd = m.thresholdInput.Update(msg)
	return m, cmd
}

func (m Model) handleSeedFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.seedInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.seedInput.Blur()
		}
	}

	var cmd tea.Cmd
	m.seedInput, cmd = m.seedInput.Update(msg)
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
	case SeedForm:
		m.seedInput.Focus()
		return m, nil
	case ThresholdForm:
		m.thresholdInput.Focus()
		return m, nil
	case DarkVariance, LightVariance, Sequence, Random:
		m.selectionMode = m.active
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
		// DarkVariance/LightVariance goes up to the row above based on active tab
		if m.focus == DarkVariance || m.focus == LightVariance {
			switch m.charControls {
			case Ascii:
				return m.setFocus(m.asciiMode)
			case Custom:
				return m.setFocus(SymbolsForm)
			}
		}
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		// From an inactive tab, navigate into the active tab's content
		downFrom := m.focus
		switch m.focus {
		case Ascii, Unicode, Custom:
			if m.charControls != m.focus {
				downFrom = m.charControls
			}
		}
		// Navigate to threshold input from DarkVar/LightVar when a variance mode is selected
		if (downFrom == DarkVariance || downFrom == LightVariance) &&
			(m.selectionMode == DarkVariance || m.selectionMode == LightVariance) {
			return m.setFocus(ThresholdForm)
		}
		// Navigate to seed input from Sequence/Random when Random is selected
		if (downFrom == Sequence || downFrom == Random) && m.selectionMode == Random {
			return m.setFocus(SeedForm)
		}
		if next, hasNext := navMap[Down][downFrom]; hasNext {
			return m.setFocus(next)
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Tab):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
		// From an inactive tab, navigate into the active tab's content
		downFrom := m.focus
		switch m.focus {
		case Ascii, Unicode, Custom:
			if m.charControls != m.focus {
				downFrom = m.charControls
			}
		}
		if next, hasNext := navMap[Down][downFrom]; hasNext {
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
