package destination

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

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

	filepath string

	ShouldClose   bool
	ShouldUnfocus bool

	IsActive bool

	width int
}

func New(w int) Model {
	filepath, _ := os.Getwd()

	return Model{
		focus: Input,

		Browser: browser.New(nil, w-2),

		filepath: filepath,

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

// View draws a control panel like this:
// |Single File   Directory
//
//	Source <path/to/file_or_dir>
//	 Select a file / directory (display file browser if sourceFilepath activated)
//	 <dir>
//	 <dir>
//	 <...>
//	Include Sub-Directories |Y  N (display if 'Directory' selected above)
//	Destination <path/to/filepath/dir>
//	 Select a directory (display file browser if filepath activated)
//	 <dir>
//	 <dir>
//	 <...>
//
// Export
func (m Model) View() string {
	return m.drawInput()
}