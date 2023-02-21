package controls

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Controls struct {
	Width, Height int
	style         lipgloss.Style
}

func New(w, h int, s lipgloss.Style) *Controls {
	return &Controls{Width: w, Height: h, style: s}
}

func (c *Controls) Init() tea.Cmd {
	return nil
}

func (c *Controls) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *Controls) View() string {
	controlStyle := c.style.Copy().Width(c.Width).Height(c.Height)
	return controlStyle.Render(fmt.Sprintf("Controls %dx%d", c.Width, c.Height))
}
