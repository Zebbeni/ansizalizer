package button

import tea "github.com/charmbracelet/bubbletea"

// Should style live on the button itself or be applied by the rendering parent?
// also, should we have separate states for selected vs active vs hovered? Or a couple of those?

type SelectCallback func()

type Button struct {
	name  string
	state State

	onSelect SelectCallback
}

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *Button) View() string {
	return b.name
}
