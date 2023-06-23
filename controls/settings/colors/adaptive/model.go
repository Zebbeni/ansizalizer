package adaptive

import (
	"image/color"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/textinput"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
)

type State int

const (
	CountForm State = iota
	IterForm
	Generate
	Save
)

type Model struct {
	focus  State
	active State

	palette palette.Model

	countInput textinput.Model
	iterInput  textinput.Model

	width, height int

	ShouldClose   bool
	ShouldUnfocus bool
	IsActive      bool
}

func New(w int) Model {
	return Model{
		focus: CountForm,

		countInput: newInput(CountForm),
		iterInput:  newInput(IterForm),

		ShouldUnfocus: false,
		IsActive:      false,

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.active {
	case CountForm:
		if m.countInput.Focused() {
			m.countInput, cmd = m.countInput.Update(msg)
			if m.countInput.Focused() == false {
				return m, event.StartAdaptingCmd
			}
			return m, cmd
		}
	case IterForm:
		if m.iterInput.Focused() {
			m.iterInput, cmd = m.iterInput.Update(msg)
			if m.iterInput.Focused() == false {
				return m, event.StartAdaptingCmd
			}
			return m, cmd
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
	inputs := m.drawInputs()
	generate := m.drawGenerateButton()
	if len(m.palette.Colors()) == 0 {
		return lipgloss.JoinVertical(lipgloss.Top, inputs, generate)
	}

	palette := lipgloss.NewStyle().Padding(0, 1, 0, 1).Render(m.palette.View())
	saveButton := m.drawSaveButton()
	content := lipgloss.JoinVertical(lipgloss.Top, inputs, palette, generate, saveButton)
	return content
}

func (m Model) Info() (int, int) {
	var count, iterations int
	count, _ = strconv.Atoi(m.countInput.Value())
	iterations, _ = strconv.Atoi(m.iterInput.Value())
	return count, iterations
}

func (m Model) GetCurrent() palette.Model {
	return m.palette
}

func (m Model) SetPalette(colors color.Palette, name string) Model {
	m.palette = palette.New(name, colors, m.width-4, 3)
	return m
}
