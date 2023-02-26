package controls

import (
	"github.com/Zebbeni/ansizalizer/component/browser"
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"

	"github.com/Zebbeni/ansizalizer/io"
)

type Model struct {
	// content is a chronological list of Models to display in our panel,
	// initialized with a main menu of options
	// The last 'content' in this list is the one we'll render in View().
	// Contents receive a callback they can use to append additional models as
	// content, (ex. a file browser may display a subdirectory as a submenu)
	// and we can go 'back' to the previous content by removing the last one.
	content []Content

	// Controls has access to state pointers, so it can affect global state
	// (when setting the app's current file or directory, for example)
	navState *browser.State

	// rendering variables
	w, h int
}

func New(nav *browser.State) Model {
	c := Model{navState: nav}
	c.content = []Content{c.BuildMainMenu()}

	return c
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, io.KeyMap.Back):
			m.removeContent()
		}
	}
	_, cmd := m.currentContent().Update(msg)
	return m, cmd
}

func (m Model) currentContent() Content {
	return m.content[len(m.content)-1]
}

func (m Model) View() string {
	content := m.currentContent().View()

	contentHeight := lipgloss.Height(content)
	_, yPosition := m.currentContent().GetActivePosition()

	activeLine := int(math.Ceil(yPosition * float64(contentHeight)))
	// get the offset needed to center the active line in the viewport.
	// (subtract 2 from the height to compensate for the viewport border)
	yOffset := activeLine - ((m.h) / 2)
	yOffset = max(0, yOffset)

	vp := viewport.New(m.w, m.h)
	vp.SetContent(content)
	vp.SetYOffset(yOffset)

	vp.Style = style.ControlsBorder.Copy().Width(m.w).Height(m.h)
	vp.Style.GetVerticalBorderSize()

	return vp.View()
}

func (m Model) Resize(w, h int) {
	m.w, m.h = w, h
}

func max(a, b int) int {
	if a > b {
		return b
	}
	return b
}
