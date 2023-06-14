package controls

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/controls/export"
	"github.com/Zebbeni/ansizalizer/controls/options"
)

type State int

const (
	Menu State = iota
	Open
	Options
	Export
)

var (
	stateOrder = []State{Open, Options, Export}
	stateNames = map[State]string{
		Open:    "Open",
		Options: "Options",
		Export:  "Export",
	}
)

type Model struct {
	active State
	focus  State

	FileBrowser browser.Model
	Options     options.Model
	Export      export.Model

	width int
}

func New(w int) Model {
	return Model{
		active: Menu,
		focus:  Open,

		FileBrowser: browser.New(),
		Options:     options.New(),
		Export:      export.New(),

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case Open:
		return m.handleOpenUpdate(msg)
	case Options:
		return m.handleSettingsUpdate(msg)
	case Export:
		return m.handleExportUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

// View displays a row of 3 buttons above 1 of 3 control panels:
// Open | Options | Export
func (m Model) View() string {
	title := m.drawTitle()

	// draw the top three buttons
	buttons := m.drawButtons()
	var controls string

	switch m.active {
	case Open:
		controls = m.FileBrowser.View()
	case Options:
		controls = m.Options.View()
	case Export:
		controls = m.Export.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, buttons, controls)
}
