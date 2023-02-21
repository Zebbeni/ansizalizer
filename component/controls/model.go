package controls

import (
	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/menu"
	"github.com/charmbracelet/bubbles/help"
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

func New(w, h int, s lipgloss.Style, keyMap help.KeyMap) *Controls {
	menu := menu.New([]item.Model{
		item.New("Open", func() {}),
		item.New("Settings", func() {}),
		item.New("Process", func() {}),
	})
	return &Controls{Width: w, Height: h, style: s, menu: menu}
}

func (c *Controls) Init() tea.Cmd {
	return nil
}

func (c *Controls) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *Controls) View() string {
	return c.menu.View()
}
