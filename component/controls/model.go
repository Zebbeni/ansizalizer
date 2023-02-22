package controls

import (
	"github.com/Zebbeni/ansizalizer/component/keyboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Controls renders the left-hand side of the interface, initially populated
// with a main menu. User actions may cause this menu to be swapped with a new
// model. Controls generally deals with managing whatever Content it currently
// holds, forwarding key and mouse messages to it, etc.
type Controls struct {
	Width, Height int
	keymap        *keyboard.Map
	style         lipgloss.Style
	content       Content

	// probably need to app state here too
}

func New(w, h int, s lipgloss.Style, k *keyboard.Map) *Controls {
	m := NewMainMenu(k)
	return &Controls{Width: w, Height: h, keymap: k, style: s, content: m}
}

func (c *Controls) Init() tea.Cmd {
	return nil
}

func (c *Controls) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *Controls) View() string {
	return c.content.View()
}

func (c *Controls) HandleKeyMsg(msg tea.KeyMsg) bool {
	return c.content.HandleKeyMsg(msg)
}
