package menu

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
)

type Model struct {
	activeItem int
	items      []item.Model
	keymap     *keyboard.Map
}

func New(b []item.Model, k *keyboard.Map) *Model {
	b[0].SetActive(true)
	return &Model{
		activeItem: 0,
		items:      b,
		keymap:     k,
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

func (m *Model) HandleKeyMsg(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, m.keymap.Nav):
		m.items[m.activeItem].SetActive(false)
		if msg.String() == "up" {
			m.activeItem = max(m.activeItem-1, 0)
		} else {
			m.activeItem = min(m.activeItem+1, len(m.items)-1)
		}
		m.items[m.activeItem].SetActive(true)
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
