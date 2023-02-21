package controls

import (
	"fmt"
	"github.com/Zebbeni/ansizalizer/component/menu"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Controls renders the left-hand side of the interface, initially a vertical
// menu. Each menu item, when selected, displays a new item to the right of the
// previous menu, cause the Controls area to expand rightward.
//
// In some cases this is a submenu. (e.g. ) In other cases it is a panel of controls
// related to the menu item. (e.g. 'Open' displays a panel to render and select
// file names)
type Controls struct {
	Width, Height int
	style         lipgloss.Style
	menu          menu.Model
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
