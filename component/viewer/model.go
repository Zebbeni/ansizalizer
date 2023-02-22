package viewer

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The Viewer returns a

type Model struct {
	style  lipgloss.Style
	keymap *help.KeyMap
}

func New(s lipgloss.Style) *Model {
	return &Model{style: s}
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

func (v *Model) HandleKeyMsg(msg tea.KeyMsg) bool {
	return false
}
