package viewer

import (
	"github.com/Zebbeni/ansizalizer/state"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state *state.Model
}

func New(state *state.Model) *Model {
	return &Model{state: state}
}

func (v *Model) Init() tea.Cmd {
	return nil
}

func (v *Model) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (v *Model) View() string {
	return "Viewer"
}
