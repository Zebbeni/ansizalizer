package export

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
	SrcInput
	SrcBrowser
	SubDirYes
	SubDirsNo
	DstInput
	DstBrowser
	Export
)

type Model struct {
	focus State

	doExportDirectory      bool
	doExportSubDirectories bool

	SourceBrowser      browser.Model
	DestinationBrowser browser.Model

	sourceFilepath      string
	destinationFilepath string

	ShouldClose   bool
	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	sourceBrowser := browser.New(nil, w-2)
	destinationBrowser := browser.New(nil, w-2)

	return Model{
		focus: ExpFile,

		SourceBrowser:      sourceBrowser,
		DestinationBrowser: destinationBrowser,

		doExportDirectory:      false,
		doExportSubDirectories: false,

		sourceFilepath: "blah/blah/blah",

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
	case SrcBrowser:
		return m.handleSrcBrowserUpdate(msg)
	case DstBrowser:
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
//	Destination <path/to/destinationFilepath/dir>
//	 Select a directory (display file browser if destinationFilepath activated)
//	 <dir>
//	 <dir>
//	 <...>
//
// Export
func (m Model) View() string {
	// draw file / directory options
	exportTypes := m.drawExportTypeOptions()
	source := m.drawSource()
	// draw sourceFilepath
	// draw subdirectory options
	// draw destinationFilepath
	// draw export button
	// join all vertically
	return lipgloss.JoinVertical(lipgloss.Center, exportTypes, source)
}
