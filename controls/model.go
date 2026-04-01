package controls

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/controls/export"
	"github.com/Zebbeni/ansizalizer/controls/filemenu"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	// modal import not needed — we use filemenu.Result pointer
	"github.com/Zebbeni/ansizalizer/global"
	"github.com/Zebbeni/ansizalizer/prefs"
)

type State int

const (
	Menu State = iota
	FileMenu
	Browse
	Settings
	Export

	numButtons = 3 // File, Browse, Settings
)

var (
	stateOrder = []State{FileMenu, Browse, Settings}
	stateNames = map[State]string{
		FileMenu: "File",
		Browse:   "Images",
		Settings: "Settings",
	}
)

type Model struct {
	active  State
	focus   State
	showing State // which tab's content is displayed (persists when returning to tab bar)

	FileDropdown filemenu.Model
	FileBrowser  browser.Model
	Settings     settings.Model
	Export       export.Model
	OpenModal    *filemenu.Result // set when a modal should be opened by the app
	ShouldQuit   bool             // set when user presses Esc at top level

	width int
}

func New(w int, dirs prefs.Dirs) Model {
	fileBrowser := browser.New(global.ImgExtensions, w)
	if dirs.Browse != "" {
		fileBrowser = browser.NewAtDir(dirs.Browse, global.ImgExtensions, w)
	}

	return Model{
		active: Menu,
		focus:  FileMenu,

		FileDropdown: filemenu.New(w),
		FileBrowser:  fileBrowser,
		Settings:     settings.New(w, dirs),
		Export:       export.New(w, dirs),

		width: w,
	}
}

func (m *Model) RefreshAllStyles() {
	m.FileBrowser.RefreshStyles()
	m.Settings.RefreshAllStyles()
}

func (m Model) DebugState() string {
	stateNames := map[State]string{Menu: "Menu", FileMenu: "FileMenu", Browse: "Browse", Settings: "Settings", Export: "Export"}
	return fmt.Sprintf("Controls: active=%s focus=%s showing=%s", stateNames[m.active], stateNames[m.focus], stateNames[m.showing])
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// File dropdown takes priority when open
	if m.FileDropdown.Open {
		return m.handleFileMenuUpdate(msg)
	}

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

// ButtonHeight returns the rendered height of the button row.
func (m Model) ButtonHeight() int {
	return lipgloss.Height(m.drawButtons())
}

// View displays a row of buttons above the active control panel.
func (m Model) View() string {
	buttons := m.drawButtons()
	var controls string

	contentState := m.showing
	if m.active != Menu && m.active != FileMenu {
		contentState = m.active
	}

	switch contentState {
	case Browse:
		browserTitle := m.drawBrowserTitle()
		controls = lipgloss.JoinVertical(lipgloss.Left, browserTitle, m.FileBrowser.View())
	case Settings:
		controls = m.Settings.View()
	case Export:
		controls = m.Export.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top, buttons, controls)
}
