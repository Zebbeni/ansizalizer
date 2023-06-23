package size

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
)

type State int
type Mode int

const (
	Fit Mode = iota
	Stretch
)

const (
	FitButton State = iota
	StretchButton
	WidthForm
	HeightForm
)

type Model struct {
	focus  State
	active State
	mode   Mode

	widthInput  textinput.Model
	heightInput textinput.Model

	ShouldUnfocus bool
	ShouldClose   bool
	IsActive      bool
}

func New() Model {
	return Model{
		focus:       FitButton,
		active:      FitButton,
		mode:        Fit,
		widthInput:  newInput(WidthForm, 50),
		heightInput: newInput(HeightForm, 40),

		ShouldUnfocus: false,
		ShouldClose:   false,
		IsActive:      false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case WidthForm:
		if m.widthInput.Focused() {
			return m.handleWidthUpdate(msg)
		}
	case HeightForm:
		if m.heightInput.Focused() {
			return m.handleHeightUpdate(msg)
		}
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
	buttonRow := m.drawButtons()
	forms := m.drawInputs()
	return lipgloss.JoinVertical(lipgloss.Left, buttonRow, forms)
}

func (m Model) Info() (Mode, int, int) {
	var width, height int
	width, _ = strconv.Atoi(m.widthInput.Value())
	height, _ = strconv.Atoi(m.heightInput.Value())
	return m.mode, width, height
}
