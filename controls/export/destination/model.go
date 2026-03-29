package destination

import (
	"os"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
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

func New(w int, dir string) Model {
	selectedDir, _ := os.Getwd()
	browserModel := browser.New(nil, w-2)
	if dir != "" {
		selectedDir = dir
		browserModel = browser.NewAtDir(dir, nil, w-2)
	}

	return Model{
		focus: Input,

		Browser: browserModel,

		selectedDir: selectedDir,

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

	selected := style.BgStyle().PaddingTop(1).Render(m.drawSelected())
	content = append(content, selected)

	if m.focus == Browser {
		content = append(content, m.drawBrowserTitle())
		content = append(content, m.Browser.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m Model) GetSelected() string {
	return m.selectedDir
}
