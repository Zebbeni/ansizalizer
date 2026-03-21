package adjust

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Up Direction = iota
	Down
)

var navMap = map[Direction]map[State]State{
	Down: {BrightnessForm: BrightnessSlider, BrightnessSlider: ContrastForm, ContrastForm: ContrastSlider},
	Up:   {ContrastSlider: ContrastForm, ContrastForm: BrightnessSlider, BrightnessSlider: BrightnessForm},
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	if m.active == m.focus {
		switch m.active {
		case BrightnessForm:
			m.brightnessInput.Blur()
			m.active = None
		case ContrastForm:
			m.contrastInput.Blur()
			m.active = None
		}
		return m, event.StartRenderToViewCmd
	}

	m.active = m.focus
	switch m.active {
	case BrightnessForm:
		m.brightnessInput.Focus()
	case ContrastForm:
		m.contrastInput.Focus()
	case BrightnessSlider, ContrastSlider:
		m.active = None
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleBrightnessUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.brightnessInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.brightnessInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.brightnessInput, cmd = m.brightnessInput.Update(msg)
	return m, cmd
}

func (m Model) handleContrastUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.contrastInput.Blur()
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.contrastInput.Blur()
		}
	}
	var cmd tea.Cmd
	m.contrastInput, cmd = m.contrastInput.Update(msg)
	return m, cmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.active {
	case BrightnessForm:
		m.brightnessInput.Blur()
		m.active = None
	case ContrastForm:
		m.contrastInput.Blur()
		m.active = None
	}

	// Handle left/right on sliders
	if m.focus == BrightnessSlider || m.focus == ContrastSlider {
		switch {
		case key.Matches(msg, event.KeyMap.Left):
			return m.handleSliderAdjust(-1)
		case key.Matches(msg, event.KeyMap.Right):
			return m.handleSliderAdjust(1)
		}
	}

	switch {
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
	return m, nil
}

func (m Model) handleSliderAdjust(delta int) (Model, tea.Cmd) {
	switch m.focus {
	case BrightnessSlider:
		val := m.Brightness() + delta
		if val < -100 {
			val = -100
		}
		if val > 100 {
			val = 100
		}
		m.brightnessInput.SetValue(strconv.Itoa(val))
	case ContrastSlider:
		val := m.Contrast() + delta
		if val < -100 {
			val = -100
		}
		if val > 100 {
			val = 100
		}
		m.contrastInput.SetValue(strconv.Itoa(val))
	}
	return m, event.StartRenderToViewCmd
}
