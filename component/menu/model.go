package menu

import (
	"github.com/Zebbeni/ansizalizer/component/item"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	activeItem int
	items      []item.Model
}

func New(b []item.Model) Model {
	b[0].SetActive(true)
	return Model{
		activeItem: 0,
		items:      b,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	items := make([]string, len(m.items))
	for i, item := range m.items {
		items[i] = item.View()
	}
	content := lipgloss.JoinVertical(lipgloss.Top, items...)
	return content
}
