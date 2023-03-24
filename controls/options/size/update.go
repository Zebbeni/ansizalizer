package size

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/io"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {FitButton: StretchButton, WidthForm: HeightForm},
	Left:  {StretchButton: FitButton, HeightForm: WidthForm},
	Up:    {WidthForm: FitButton, HeightForm: StretchButton},
	Down:  {FitButton: WidthForm, StretchButton: HeightForm},
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	if m.active == m.focus && (m.active == FitButton || m.active == StretchButton) {
		m.ShouldClose = true
		return m, nil
	}

	m.active = m.focus
	switch m.active {
	case FitButton:
		m.mode = Fit
	case StretchButton:
		m.mode = Stretch
	case WidthForm:
		m.widthInput.Focus()
	case HeightForm:
		m.heightInput.Focus()
	}
	return m, io.StartRenderCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, io.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, io.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, io.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, io.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			m.focus = next
		}
	}

	return m, cmd
}
