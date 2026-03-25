package adjust

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	BrightnessForm State = iota
	BrightnessSlider
	ContrastForm
	ContrastSlider
	None
)

type Model struct {
	focus  State
	active State

	brightnessInput textinput.Model
	contrastInput   textinput.Model

	ShouldClose bool
	IsActive    bool
	width       int
}

func New(w int) Model {
	return Model{
		focus:           BrightnessForm,
		active:          None,
		brightnessInput: newInput("Brightness", 0),
		contrastInput:   newInput("Contrast", 0),
		ShouldClose:     false,
		IsActive:        false,
		width:           w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd1, cmd2 tea.Cmd
	newM := m

	switch m.active {
	case BrightnessForm:
		if m.brightnessInput.Focused() {
			newM, cmd1 = newM.handleBrightnessUpdate(msg)
		}
	case ContrastForm:
		if m.contrastInput.Focused() {
			newM, cmd1 = newM.handleContrastUpdate(msg)
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
	brightness := m.drawBrightnessForm()
	brightnessSlider := m.drawSlider(BrightnessSlider, m.Brightness())
	contrast := m.drawContrastForm()
	contrastSlider := m.drawSlider(ContrastSlider, m.Contrast())
	return lipgloss.JoinVertical(lipgloss.Left, brightness, brightnessSlider, contrast, contrastSlider)
}

func (m Model) Brightness() int {
	return clampInput(m.brightnessInput.Value())
}

func (m Model) Contrast() int {
	return clampInput(m.contrastInput.Value())
}

func (m *Model) SetConfig(brightness, contrast int) {
	m.brightnessInput.SetValue(strconv.Itoa(brightness))
	m.contrastInput.SetValue(strconv.Itoa(contrast))
}

func (m *Model) ResetFocus() {
	m.focus = BrightnessForm
	m.active = None
	m.brightnessInput.Blur()
	m.contrastInput.Blur()
}

func (m Model) Summary() string {
	return "Bright: " + m.brightnessInput.Value() + "\nContrast: " + m.contrastInput.Value()
}

func clampInput(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	if val < -100 {
		return -100
	}
	if val > 100 {
		return 100
	}
	return val
}
