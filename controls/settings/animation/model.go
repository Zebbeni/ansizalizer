package animation

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	DelayForm State = iota
	None
)

type Model struct {
	focus  State
	active State

	delayInput numberinput.Model

	ShouldClose   bool
	IsActive      bool
	AnimatedImage bool
}

func New() Model {
	return Model{
		focus:         DelayForm,
		active:        None,
		delayInput:    newInput("Delay (ms)", 100),
		ShouldClose:   false,
		IsActive:      false,
		AnimatedImage: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd1, cmd2 tea.Cmd
	newM := m

	switch m.active {
	case DelayForm:
		if m.delayInput.Focused() {
			newM, cmd1 = newM.handleDelayUpdate(msg)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			newM, cmd2 = newM.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			newM, cmd2 = newM.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Esc):
			newM, cmd2 = newM.handleEsc()
		}
	}
	return newM, tea.Batch(cmd1, cmd2)
}

func (m Model) View() string {
	return m.drawDelayForm()
}

func (m Model) Delay() time.Duration {
	return time.Duration(m.delayInput.IntValue()) * time.Millisecond
}

func (m Model) Summary() string {
	return "Delay: " + m.delayInput.Value() + "ms"
}

func (m Model) DelayMs() int {
	return m.delayInput.IntValue()
}

func (m *Model) SetDelayMs(ms int) {
	m.delayInput.SetIntValue(ms)
}

func (m *Model) ResetFocus() {
	m.focus = DelayForm
	m.active = None
	m.delayInput.Blur()
}
