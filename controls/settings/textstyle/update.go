package textstyle

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
		BoldOn: BoldOff, ItalicOn: ItalicOff,
		UnderlineOn: UnderlineOff, StrikethroughOn: StrikethroughOff,
	},
	Left: {
		BoldOff: BoldOn, ItalicOff: ItalicOn,
		UnderlineOff: UnderlineOn, StrikethroughOff: StrikethroughOn,
	},
	Down: {
		BoldOn: ItalicOn, BoldOff: ItalicOff,
		ItalicOn: UnderlineOn, ItalicOff: UnderlineOff,
		UnderlineOn: StrikethroughOn, UnderlineOff: StrikethroughOff,
	},
	Up: {
		ItalicOn: BoldOn, ItalicOff: BoldOff,
		UnderlineOn: ItalicOn, UnderlineOff: ItalicOff,
		StrikethroughOn: UnderlineOn, StrikethroughOff: UnderlineOff,
	},
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, ok := navMap[Right][m.focus]; ok {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, ok := navMap[Left][m.focus]; ok {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, ok := navMap[Down][m.focus]; ok {
			m.focus = next
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Tab):
		if next, ok := navMap[Right][m.focus]; ok {
			m.focus = next
		} else if next, ok := navMap[Down][m.focus]; ok {
			m.focus = next
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, ok := navMap[Up][m.focus]; ok {
			m.focus = next
		} else {
			m.IsActive = false
			m.ShouldClose = true
		}
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case BoldOn, BoldOff:
		m.bold = !m.bold
	case ItalicOn, ItalicOff:
		m.italic = !m.italic
	case UnderlineOn, UnderlineOff:
		m.underline = !m.underline
	case StrikethroughOn, StrikethroughOff:
		m.strikethrough = !m.strikethrough
	}
	return m, event.StartRenderToViewCmd
}
