package menu

import (
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
)

type Model struct {
	activeItem int
	Items      []item.Model
	keymap     *keyboard.Map
}

func New(b []item.Model, k *keyboard.Map) *Model {
	return &Model{
		activeItem: 0,
		Items:      b,
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
	items := make([]string, len(m.Items))
	for i, item := range m.Items {
		items[i] = item.View()
		if m.activeItem == i {
			items[i] = style.ActiveItem.Render(items[i])
		}
	}
	content := lipgloss.JoinVertical(lipgloss.Top, items...)
	return content
}

func (m *Model) HandleKeyMsg(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, m.keymap.Nav):
		if msg.String() == "up" {
			m.activeItem = max(m.activeItem-1, 0)
		} else {
			m.activeItem = min(m.activeItem+1, len(m.Items)-1)
		}
		return true
	case key.Matches(msg, m.keymap.Enter):
		m.Items[m.activeItem].OnSelect()
		return true
	}

	return false
}

func (m *Model) GetActivePosition() (float64, float64) {
	if len(m.Items) == 0 {
		return 0, 0
	}
	return 0, float64(m.activeItem) / float64(len(m.Items))
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
