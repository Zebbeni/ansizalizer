package textstyle

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
	case BoldOn:
		m.bold = true
	case BoldOff:
		m.bold = false
	case ItalicOn:
		m.italic = true
	case ItalicOff:
		m.italic = false
	case UnderlineOn:
		m.underline = true
	case UnderlineOff:
		m.underline = false
	case StrikethroughOn:
		m.strikethrough = true
	case StrikethroughOff:
		m.strikethrough = false
	}
	return m, event.StartRenderToViewCmd
}
