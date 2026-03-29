package sampling

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func (m *Model) SetFunctionByName(name string) {
	for fn, n := range NameMap {
		if n == name {
			m.Function = fn
			return
		}
	}
}

func (m *Model) SetListActive(active bool) {
	m.list.SetDelegate(NewDelegate(active))
}

func (m Model) FunctionName() string {
	if name, ok := NameMap[m.Function]; ok {
		return name
	}
	return "Unknown"
}

func (m Model) View() string {
	prompt := style.DimmedTitle.Copy().Render("Select Method")
	menu := m.list.View()
	content := lipgloss.JoinVertical(lipgloss.Left, prompt, menu)
	return content
}
