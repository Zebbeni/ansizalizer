package viewer

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/event"
)

type Model struct {
	imgString string
	settings  settings.Model

	WaitingOnRender bool
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case event.FinishRenderToViewMsg:
		return m.handleFinishRenderMsg(msg)
	}
	return m, nil
}

func (m Model) View() string {
	if m.WaitingOnRender {
		return ""
	}
	return m.imgString
}
