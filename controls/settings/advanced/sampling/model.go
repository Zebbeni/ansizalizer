package sampling

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	Function resize.InterpolationFunction

	list list.Model

	IsActive    bool
	ShouldClose bool
}

func New(w int) Model {
	items := menuItems()
	selected := items[0].(item)
	menu := newMenu(items, w, len(items))

	return Model{
		Function:    selected.Function,
		list:        menu,
		IsActive:    false,
		ShouldClose: false,
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
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	prompt := style.DimmedTitle.Copy().Render("Select Method")
	menu := m.list.View()
	content := lipgloss.JoinVertical(lipgloss.Left, prompt, menu)
	return lipgloss.NewStyle().Padding(0, 1).Render(content)
}
