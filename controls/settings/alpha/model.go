package alpha

import (
	"strconv"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Input State = iota
	AlphaYes
	AlphaNo
	UseAlpha
	TrimAlpha
	TrimAlphaYes
	TrimAlphaNo
	ThresholdForm
)

type Model struct {
	focus          State
	useAlpha       bool
	trimAlpha      bool
	thresholdInput textinput.Model
	ShouldUnfocus  bool
	IsActive       bool
	width          int
	AlphaImage     bool
}

func New(w int) Model {
	input := textinput.New()
	input.Prompt = "Threshold "
	input.CharLimit = 5
	input.SetValue("0.0")

	return Model{
		focus:          AlphaYes,
		useAlpha:       true,
		trimAlpha:      false,
		thresholdInput: input,
		IsActive:       false,
		width:          w,
		AlphaImage:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == ThresholdForm && m.thresholdInput.Focused() {
		return m.handleThresholdUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, nil
}

func (m Model) handleThresholdUpdate(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) View() string {
	content := make([]string, 0, 5)
	content = append(content, m.drawAlphaOptions())
	content = append(content, m.drawAlphaTrimOptions())
	if m.useAlpha {
		content = append(content, m.drawThresholdInput())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) ShouldOutputAlpha() bool {
	return m.useAlpha
}

func (m Model) TrimAlpha() bool {
	return m.trimAlpha
}

func (m Model) AlphaThreshold() float64 {
	val, err := strconv.ParseFloat(m.thresholdInput.Value(), 64)
	if err != nil || val < 0 {
		return 0
	}
	return val
}

func (m *Model) SetAlphaThreshold(t float64) {
	m.thresholdInput.SetValue(strconv.FormatFloat(t, 'f', 2, 64))
}

func (m *Model) SetConfig(useAlpha, trimAlpha bool) {
	m.useAlpha = useAlpha
	m.trimAlpha = trimAlpha
}

func (m *Model) ResetFocus() {
	m.focus = AlphaYes
}

func (m Model) Summary() string {
	if !m.useAlpha {
		return "Alpha: Off"
	}
	s := "Alpha: On"
	if m.trimAlpha {
		s += " | Trim: Yes"
	}
	if t := m.AlphaThreshold(); t > 0 {
		s += " | Thresh: " + m.thresholdInput.Value()
	}
	return s
}
