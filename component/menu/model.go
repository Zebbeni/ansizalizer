package menu

import (
	"github.com/Zebbeni/ansizalizer/component/button"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	activeIdx int
	buttons   []button.Button
}

func New(b []button.Button) *Model {
	return &Model{
		activeIdx: 0,
		buttons:   b,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return ""
}
