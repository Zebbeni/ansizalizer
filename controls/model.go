package controls

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/controls/export"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/global"
)

type State int

const (
	Menu State = iota
	Browse
	Settings
	Export

	numButtons = 3
)

var (
	stateOrder = []State{Browse, Settings, Export}
	stateNames = map[State]string{
		Browse:   "Browse",
		Settings: "Settings",
		Export:   "Export",
	}
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
		focus:  Browse,

		FileBrowser: browser.New(global.ImgExtensions, w),
		Settings:    settings.New(w),
		Export:      export.New(w),

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case Browse:
		return m.handleOpenUpdate(msg)
	case Settings:
		return m.handleSettingsUpdate(msg)
	case Export:
		return m.handleExportUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

// View displays a row of 3 buttons above 1 of 3 control panels:
// Browse | Settings | Export
func (m Model) View() string {
	title := m.drawTitle()

	// draw the top three buttons
	buttons := m.drawButtons()
	var controls string

	switch m.active {
	case Browse:
		browserTitle := m.drawBrowserTitle()
		controls = lipgloss.JoinVertical(lipgloss.Left, browserTitle, m.FileBrowser.View())
	case Settings:
		controls = m.Settings.View()
	case Export:
		controls = m.Export.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, buttons, controls)
}
