package controls

import (
	"github.com/Zebbeni/ansizalizer/component/keyboard"
	"github.com/Zebbeni/ansizalizer/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Controls renders the left-hand side of the interface, initially populated
// with a main menu. User actions may cause this menu to be swapped with a new
// model. Controls generally deals with managing whatever Content it currently
// holds, forwarding key and mouse messages to it, etc.
type Controls struct {
	state  *state.Model
	keymap *keyboard.Map

	// content is a chronological list of Content objects displayed in the
	// Controls panel. The last item listed is the one to display.
	//
	// When new content is added, (e.g. When 'Open' adds a file browser)
	// we append the new Content to the content stack. When the user presses
	// 'esc' to go back, we remove the browser from the Content stack.
	content []Content
}

func New(s *state.Model, k *keyboard.Map) *Controls {
	c := &Controls{state: s, keymap: k}
	c.content = []Content{NewMainMenu(s, k, c.addContent)}

	return c
}

func (c *Controls) Init() tea.Cmd {
	return nil
}

func (c *Controls) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *Controls) View() string {
	return c.content[len(c.content)-1].View()
}

func (c *Controls) HandleKeyMsg(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, c.keymap.Back):
		c.removeContent()
		return true
	}
	return c.content[len(c.content)-1].HandleKeyMsg(msg)
}

func (c *Controls) addContent(content Content) {
	c.content = append(c.content, content)
}

func (c *Controls) removeContent() {
	if len(c.content) > 1 {
		c.content = c.content[:len(c.content)-1]
	}
}
