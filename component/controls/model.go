package controls

import tea "github.com/charmbracelet/bubbletea"

type Controls struct {
}

func New() *Controls {
	return &Controls{}
}

func (c *Controls) Init() tea.Cmd {
	return nil
}

func (c *Controls) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *Controls) View() string {
	return "Controls"
}
