package destination

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Input State = iota
	Browser
)

type Model struct {
	focus State

	Browser browser.Model

	selectedDir string

	ShouldClose   bool
	ShouldUnfocus bool

	IsActive bool

	width int
}

func New(w int) Model {
	filepath, _ := os.Getwd()

	return Model{
		focus: Input,

		Browser: browser.New(nil, drawBrowserTitle(), w-2),

		selectedDir: filepath,

		width:       w,
		ShouldClose: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focus {
	case Browser:
		return m.handleDstBrowserUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	content := make([]string, 0, 5)

	selected := lipgloss.NewStyle().PaddingTop(1).Render(m.drawSelected())
	content = append(content, selected)

	if m.focus == Browser {
		content = append(content, m.Browser.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) GetSelected() string {
	return m.selectedDir
}
