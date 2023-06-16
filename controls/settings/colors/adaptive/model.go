package adaptive

import (
	"image/color"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/textinput"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors/description"
	"github.com/Zebbeni/ansizalizer/io"
)

type State int

const (
	CountForm State = iota
	IterForm
	Generate
)

type Model struct {
	focus  State
	active State

	Palette color.Palette

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
				return m, io.BuildAdaptingCmd()
			}
			return m, cmd
		}
	case IterForm:
		if m.iterInput.Focused() {
			m.iterInput, cmd = m.iterInput.Update(msg)
			if m.iterInput.Focused() == false {
				return m, io.BuildAdaptingCmd()
			}
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, io.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, io.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, io.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	inputs := m.drawInputs()
	generate := m.drawGenerateButton()
	if len(m.Palette) == 0 {
		return lipgloss.JoinVertical(lipgloss.Top, inputs, generate)
	}

	palette := lipgloss.NewStyle().Padding(0, 1, 0, 1).Render(description.Palette(m.Palette, 25, 10))
	return lipgloss.JoinVertical(lipgloss.Top, inputs, generate, palette)
}

func (m Model) Info() (int, int) {
	var count, iterations int
	count, _ = strconv.Atoi(m.countInput.Value())
	iterations, _ = strconv.Atoi(m.iterInput.Value())
	return count, iterations
}

func (m Model) GetCurrent() color.Palette {
	return m.Palette
}
