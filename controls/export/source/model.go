package source

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	ExpFile State = iota
	ExpDirectory
	Input
	Browser
	SubDirYes
	SubDirsNo
)

type Model struct {
	focus State

	doExportDirectory      bool
	doExportSubDirectories bool

	Browser      browser.Model
	selectedDir  string
	selectedFile string

	ShouldClose   bool
	ShouldUnfocus bool

	IsActive bool

	width int
}

func New(w int) Model {
	browserModel := browser.New(nil, w-2)

	return Model{
		focus: ExpDirectory,

		Browser: browserModel,

		doExportDirectory:      false,
		doExportSubDirectories: false,

		selectedDir:  "",
		selectedFile: "",

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
		return m.handleSrcBrowserUpdate(msg)
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
//	 Select a file / directory (display file browser if filepath activated)
//	 <dir>
//	 <dir>
//	 <...>
//	Include Sub-Directories |Y  N (display if 'Directory' selected above)
//	Destination <path/to/destinationFilepath/dir>
//	 Select a directory (display file browser if destinationFilepath activated)
//	 <dir>
//	 <dir>
//	 <...>
//
// Export
func (m Model) View() string {
	content := make([]string, 0, 5)
	content = append(content, m.drawExportTypeOptions())

	selected := lipgloss.NewStyle().PaddingTop(1).Render(m.drawSelected())
	content = append(content, selected)

	if m.focus == Browser {
		content = append(content, m.Browser.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}
