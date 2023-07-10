package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	msg   string
	width int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case event.DisplayMsg:
		m.msg = string(msg)
	}
	return m, nil
}

func (m Model) View() string {
	// TODO: Switch style based on event type (warning, info, etc.)
	displayStyle := style.ExtraDimTitle.Copy().Width(m.width - 2)
	return displayStyle.Border(lipgloss.RoundedBorder()).BorderForeground(style.ExtraDimColor).Render(m.msg)
}

func (m Model) SetWidth(w int) Model {
	m.width = w
	return m
}
