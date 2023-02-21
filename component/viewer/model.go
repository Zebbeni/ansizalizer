package viewer

import tea "github.com/charmbracelet/bubbletea"

type Viewer struct {
}

func New() *Viewer {
	return &Viewer{}
}

func (c *Viewer) Init() tea.Cmd {
	return nil
}

func (c *Viewer) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *Viewer) View() string {
	return "Viewer"
}
