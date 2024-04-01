package size

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
	Right: {FitButton: StretchButton, WidthForm: HeightForm},
	Left:  {StretchButton: FitButton, HeightForm: WidthForm},
	Up:    {WidthForm: FitButton, HeightForm: StretchButton, CharRatioForm: HeightForm},
	Down:  {FitButton: WidthForm, StretchButton: HeightForm, WidthForm: CharRatioForm, HeightForm: CharRatioForm},
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
	case CharRatioForm:
		m.charRatioInput.Focus()
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleWidthUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.widthInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.widthInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.widthInput, cmd = m.widthInput.Update(msg)
	return m, cmd
}

func (m Model) handleHeightUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.heightInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.heightInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.heightInput, cmd = m.heightInput.Update(msg)
	return m, cmd
}

func (m Model) handleCharRatioUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.charRatioInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.charRatioInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.charRatioInput, cmd = m.charRatioInput.Update(msg)
	return m, cmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			m.focus = next
		} else {
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			m.focus = next
		} else {
			m.ShouldClose = true
		}
	}

	return m, cmd
}
