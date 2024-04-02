package size

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
)

const DEFAULT_CHAR_W_TO_H_RATIO = 0.5

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
	CharRatioForm
	None
)

type Model struct {
	focus  State
	active State
	mode   Mode

	widthInput     textinput.Model
	heightInput    textinput.Model
	charRatioInput textinput.Model

	ShouldUnfocus bool
	ShouldClose   bool
	IsActive      bool
}

func New() Model {
	return Model{
		focus:          FitButton,
		active:         None,
		mode:           Fit,
		widthInput:     newInput(WidthForm, 50),
		heightInput:    newInput(HeightForm, 40),
		charRatioInput: newFloatInput(CharRatioForm, DEFAULT_CHAR_W_TO_H_RATIO),

		ShouldUnfocus: false,
		ShouldClose:   false,
		IsActive:      false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd1, cmd2 tea.Cmd
	newM := m

	switch m.active {
	case WidthForm:
		if m.widthInput.Focused() {
			newM, cmd1 = newM.handleWidthUpdate(msg)
		}
	case HeightForm:
		if m.heightInput.Focused() {
			newM, cmd1 = newM.handleHeightUpdate(msg)
		}
	case CharRatioForm:
		if m.charRatioInput.Focused() {
			newM, cmd1 = newM.handleCharRatioUpdate(msg)
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
	buttonRow := m.drawButtons()
	forms := m.drawSizeForms()
	ratioForm := m.drawCharRatioForm()
	return lipgloss.JoinVertical(lipgloss.Left, buttonRow, forms, ratioForm)
}

func (m Model) Info() (Mode, int, int, float64) {
	var width, height int
	width, _ = strconv.Atoi(m.widthInput.Value())
	height, _ = strconv.Atoi(m.heightInput.Value())
	charRatio, err := strconv.ParseFloat(m.charRatioInput.Value(), 64)
	if err != nil {
		charRatio = DEFAULT_CHAR_W_TO_H_RATIO
	}
	return m.mode, width, height, charRatio
}
