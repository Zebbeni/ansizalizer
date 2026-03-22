package animation

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

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

	delayInput textinput.Model

	ShouldClose bool
	IsActive    bool
}

func New() Model {
	return Model{
		focus:       DelayForm,
		active:      None,
		delayInput:  newInput("Delay (ms)", 100),
		ShouldClose: false,
		IsActive:    false,
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
	val, err := strconv.Atoi(m.delayInput.Value())
	if err != nil || val < 10 {
		return 100 * time.Millisecond
	}
	if val > 2000 {
		return 2000 * time.Millisecond
	}
	return time.Duration(val) * time.Millisecond
}

func (m Model) Summary() string {
	return "Delay: " + m.delayInput.Value() + "ms"
}

func (m *Model) ResetFocus() {
	m.focus = DelayForm
	m.active = None
	m.delayInput.Blur()
}
