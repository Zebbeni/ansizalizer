package viewer

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Viewer struct {
	Width, Height int
	style         lipgloss.Style
}

func New(w, h int, s lipgloss.Style) *Viewer {
	return &Viewer{Width: w, Height: h, style: s}
}

func (v *Viewer) Init() tea.Cmd {
	return nil
}

func (v *Viewer) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (v *Viewer) View() string {
	style := v.style.Copy().Width(v.Width).Height(v.Height)
	return style.Render(fmt.Sprintf("Controls %dx%d", v.Width, v.Height))
}
