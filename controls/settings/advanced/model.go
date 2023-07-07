package advanced

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/sampling"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Menu State = iota
	Sampling
	Dithering
	SamplingControls
	DitheringControls
)

type Model struct {
	focus       State
	active      State
	activeTab   State
	sampling    sampling.Model
	dithering   dithering.Model
	ShouldClose bool
	IsActive    bool
	width       int
}

func New(w int) Model {
	return Model{
		focus:       Sampling,
		active:      Menu,
		activeTab:   Sampling,
		sampling:    sampling.New(w - 2),
		dithering:   dithering.New(w - 2),
		ShouldClose: false,
		IsActive:    false,
		width:       w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case SamplingControls:
		return m.handleSamplingUpdate(msg)
	case DitheringControls:
		return m.handleDitheringUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.drawTabs()
}

func (m Model) SamplingFunction() resize.InterpolationFunction {
	return m.sampling.Function
}
