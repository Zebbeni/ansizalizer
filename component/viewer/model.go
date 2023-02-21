package viewer

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width, Height int
	style         lipgloss.Style
}

func New(w, h int, s lipgloss.Style) *Model {
	return &Model{Width: w, Height: h, style: s}
}

func (v *Model) Init() tea.Cmd {
	return nil
}

func (v *Model) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (v *Model) View() string {
	style := v.style.Copy().Width(v.Width).Height(v.Height)
	return style.Render(fmt.Sprintf("Controls %dx%d", v.Width, v.Height))
}
