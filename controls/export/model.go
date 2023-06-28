package export

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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

	source      string
	destination string

	width       int
	ShouldClose bool
}

func New(w int) Model {
	return Model{
		focus: ExpFile,

		doExportDirectory:      false,
		doExportSubDirectories: false,

		width:       w,
		ShouldClose: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			m.ShouldClose = true
		}
	}
	return m, nil
}

// View draws a control panel like this:
// |Single File   Directory
//
//	Source <path/to/file_or_dir>
//	 Select a file / directory (display file browser if source activated)
//	 <dir>
//	 <dir>
//	 <...>
//	Include Sub-Directories |Y  N (display if 'Directory' selected above)
//	Destination <path/to/destination/dir>
//	 Select a directory (display file browser if destination activated)
//	 <dir>
//	 <dir>
//	 <...>
//
// Export
func (m Model) View() string {
	// draw file / directory options
	exportTypes := m.drawExportTypeOptions()
	// draw source
	// draw subdirectory options
	// draw destination
	// draw export button
	// join all vertically
	return lipgloss.JoinVertical(lipgloss.Center, exportTypes)
}
