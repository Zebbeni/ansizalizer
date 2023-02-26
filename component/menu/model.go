package menu

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/Zebbeni/ansizalizer/io"
)

type Model struct {
	activeItem int
	Items      []item.Model
}

func New(b []item.Model) Model {
	return Model{
		activeItem: 0,
		Items:      b,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, io.KeyMap.Nav):
			if msgType.String() == "up" {
				m.activeItem = max(m.activeItem-1, 0)
			} else {
				m.activeItem = min(m.activeItem+1, len(m.Items)-1)
			}
		case key.Matches(msgType, io.KeyMap.Enter):
			// should OnSelect return a command?
			m.Items[m.activeItem].OnSelect()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	items := make([]string, len(m.Items))
	for i, item := range m.Items {
		items[i] = item.View()
		if m.activeItem == i {
			items[i] = style.ActiveItem.Render(items[i])
		}
	}
	content := lipgloss.JoinVertical(lipgloss.Top, items...)
	// TODO: We're padding the bottom to prevent viewport truncation due to
	// https://github.com/charmbracelet/bubbles/issues/336. Remove when fixed.
	content = lipgloss.NewStyle().PaddingBottom(2).Render(content)

	return content
}

func (m Model) ResetActive() {
	m.activeItem = 0
}

func (m Model) GetActivePosition() (float64, float64) {
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
