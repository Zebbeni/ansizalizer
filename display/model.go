package display

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	msg string
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
	displayStyle := style.DimmedTitle.Copy().PaddingTop(1).Width(100)
	return displayStyle.Render(m.msg)
}
