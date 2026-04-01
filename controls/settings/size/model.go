package size

import (
	"strconv"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/event"
)

const DEFAULT_CHAR_W_TO_H_RATIO = 0.46

type State int
type Mode int

const (
	Fit Mode = iota
	Fill
	Stretch
)

const (
	FitButton State = iota
	FillButton
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

	widthInput     numberinput.Model
	heightInput    numberinput.Model
	charRatioInput numberinput.Model

	ShouldUnfocus bool
	ShouldClose   bool
	IsActive      bool
}

func New() Model {
	return Model{
		focus:          WidthForm,
		active:         None,
		mode:           Fit,
		widthInput:     newInput(WidthForm, 150),
		heightInput:    newInput(HeightForm, 50),
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
	switch m.active {
	case WidthForm:
		if m.widthInput.Focused() {
			return m.handleWidthUpdate(msg)
		}
	case HeightForm:
		if m.heightInput.Focused() {
			return m.handleHeightUpdate(msg)
		}
	case CharRatioForm:
		if m.charRatioInput.Focused() {
			return m.handleCharRatioUpdate(msg)
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
	forms := m.drawSizeForms()
	modeRow := m.drawModeButtons()
	ratioForm := m.drawCharRatioForm()
	return lipgloss.JoinVertical(lipgloss.Left, forms, modeRow, ratioForm)
}

func (m Model) Info() (Mode, int, int, float64) {
	width := m.widthInput.IntValue()
	height := m.heightInput.IntValue()
	charRatio := m.charRatioInput.FloatValue()
	return m.mode, width, height, charRatio
}

func (m *Model) SetConfig(mode Mode, width, height int, charRatio float64) {
	m.mode = mode
	m.widthInput.SetIntValue(width)
	m.heightInput.SetIntValue(height)
	m.charRatioInput.SetValue(strconv.FormatFloat(charRatio, 'f', -1, 64))
}

func (m *Model) ResetFocus() {
	m.focus = WidthForm
	m.active = None
	m.widthInput.Blur()
	m.heightInput.Blur()
	m.charRatioInput.Blur()
}

func (m Model) Summary() string {
	modeNames := map[Mode]string{Fit: "Fit", Fill: "Fill", Stretch: "Stretch"}
	return "Size: " + m.widthInput.Value() + "x" + m.heightInput.Value() + "\nMode: " + modeNames[m.mode]
}
