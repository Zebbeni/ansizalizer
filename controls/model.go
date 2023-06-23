package controls

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/controls/export"
	"github.com/Zebbeni/ansizalizer/controls/settings"
)

type State int

const (
	Menu State = iota
	Open
	Settings
	Export

	numButtons = 3
)

var (
	stateOrder = []State{Open, Settings, Export}
	stateNames = map[State]string{
		Open:     "Open",
		Settings: "Settings",
		Export:   "Export",
	}
	imgExtensions = []string{".jpg", ".png"}
)

type Model struct {
	active State
	focus  State

	FileBrowser browser.Model
	Settings    settings.Model
	Export      export.Model

	width int
}

func New(w int) Model {
	return Model{
		active: Menu,
		focus:  Open,

		FileBrowser: browser.New(imgExtensions, w),
		Settings:    settings.New(w),
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
	case Settings:
		return m.handleSettingsUpdate(msg)
	case Export:
		return m.handleExportUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

// View displays a row of 3 buttons above 1 of 3 control panels:
// Open | Settings | Export
func (m Model) View() string {
	title := m.drawTitle()

	// draw the top three buttons
	buttons := m.drawButtons()
	var controls string

	switch m.active {
	case Open:
		controls = m.FileBrowser.View()
	case Settings:
		controls = m.Settings.View()
	case Export:
		controls = m.Export.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, buttons, controls)
}
