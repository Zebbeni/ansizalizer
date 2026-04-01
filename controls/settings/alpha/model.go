package alpha

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Input State = iota
	UseAlpha
	TrimAlpha
	ThresholdForm
)

type Model struct {
	focus          State
	useAlpha       bool
	trimAlpha      bool
	thresholdInput numberinput.Model
	ShouldClose    bool
	IsActive       bool
	width          int
	AlphaImage     bool
}

func New(w int) Model {
	return Model{
		focus:    ThresholdForm,
		useAlpha: true,
		trimAlpha: false,
		thresholdInput: numberinput.New(numberinput.Options{
			Prompt:    "Render Threshold ",
			CharLimit: 4,
			IsFloat:   true,
			Min:       numberinput.FloatPtr(0),
			Max:       numberinput.FloatPtr(1),
			Default:   0.5,
		}),
		IsActive:   false,
		width:      w,
		AlphaImage: true,
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
	content = append(content, m.drawThresholdInput())
	content = append(content, m.drawAlphaOptions())
	if m.useAlpha {
		content = append(content, m.drawAlphaTrimOptions())
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
	return m.thresholdInput.FloatValue()
}

func (m *Model) SetAlphaThreshold(t float64) {
	m.thresholdInput.SetFloatValue(t, 2)
}

func (m *Model) SetConfig(useAlpha, trimAlpha bool) {
	m.useAlpha = useAlpha
	m.trimAlpha = trimAlpha
}

func (m *Model) ResetFocus() {
	m.focus = ThresholdForm
}

func (m Model) Summary() string {
	s := "Thresh: " + m.thresholdInput.Value()
	if m.useAlpha {
		s += " | Transparent"
		if m.trimAlpha {
			s += " | Trim"
		}
	}
	return s
}
