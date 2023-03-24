package adaptive

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/io"
)

type Model struct {
	ShouldUnfocus bool
}

func New() Model {
	return Model{
		ShouldUnfocus: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, io.KeyMap.Esc):
			m.ShouldUnfocus = true
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "Adaptive"
}
