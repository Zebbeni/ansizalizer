package alpha

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	focus         State
	useAlpha      bool
	trimAlpha     bool
	ShouldUnfocus bool
	IsActive      bool
	width         int
	AlphaImage    bool
}

func New(w int) Model {

	return Model{
		focus:      AlphaYes,
		useAlpha:   true,
		trimAlpha:  false,
		IsActive:   false,
		width:      w,
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

func (m *Model) SetConfig(useAlpha, trimAlpha bool) {
	m.useAlpha = useAlpha
	m.trimAlpha = trimAlpha
}

func (m *Model) ResetFocus() {
	m.focus = AlphaYes
}

func (m Model) Summary() string {
	if !m.useAlpha {
		return "Alpha: Off"
	}
	if m.trimAlpha {
		return "Alpha: On | Trim: Yes"
	}
	return "Alpha: On | Trim: No"
}
