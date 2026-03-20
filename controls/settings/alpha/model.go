package alpha

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Input State = iota
	AlphaYes
	AlphaNo
	UseAlpha
	TrimAlpha
	TrimAlphaYes
	TrimAlphaNo
)

type Model struct {
	focus State
	useAlpha bool
	trimAlpha bool
	ShouldUnfocus bool
	IsActive bool
	width int
	AlphaImage bool
}

func New(w int) Model {

	return Model{
		focus: AlphaYes,
		useAlpha: true,
		trimAlpha: false,
		IsActive: false,
		width: w,
		AlphaImage: true,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) View() string {
	content := make([]string, 0, 5)
	content = append(content, m.drawAlphaOptions())
	content = append(content, m.drawAlphaTrimOptions())

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) ShouldOutputAlpha() bool {
	return m.useAlpha
}

func (m Model) TrimAlpha() bool {
	return m.trimAlpha
}
