package filemenu

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Item int

const (
	LoadSettings Item = iota
	SaveSettings
	Theme
	Export
	Quit
)

type SubItem int

const (
	NoSub SubItem = iota
	ThemeLightOnTransparent
	ThemeLightOnDark
	ThemeLightOnDarkPaletted
	ThemeDarkOnTransparent
	ThemeDarkOnLight
	ThemeDarkOnLightPaletted
	ExportCurrentFile
	ExportBatchProcess
)

var menuItems = []Item{LoadSettings, SaveSettings, Theme, Export, Quit}

var menuLabels = map[Item]string{
	LoadSettings: "Load Settings",
	SaveSettings: "Save Settings",
	Theme:        "Theme         ▸",
	Export:       "Export        ▸",
	Quit:         "Quit",
}

var hasSubmenu = map[Item]bool{Theme: true, Export: true}

var themeSubItems = []SubItem{
	ThemeLightOnTransparent, ThemeLightOnDark, ThemeLightOnDarkPaletted,
	ThemeDarkOnTransparent, ThemeDarkOnLight, ThemeDarkOnLightPaletted,
}

var themeSubLabels = map[SubItem]string{
	ThemeLightOnTransparent:  "Light on Transparent",
	ThemeLightOnDark:         "Light on Dark",
	ThemeLightOnDarkPaletted: "Light on Dark (Paletted)",
	ThemeDarkOnTransparent:   "Dark on Transparent",
	ThemeDarkOnLight:         "Dark on Light",
	ThemeDarkOnLightPaletted: "Dark on Light (Paletted)",
}

var themeSubToName = map[SubItem]string{
	ThemeLightOnTransparent:  "Light on Transparent",
	ThemeLightOnDark:         "Light on Dark",
	ThemeLightOnDarkPaletted: "Light on Dark (Paletted)",
	ThemeDarkOnTransparent:   "Dark on Transparent",
	ThemeDarkOnLight:         "Dark on Light",
	ThemeDarkOnLightPaletted: "Dark on Light (Paletted)",
}

var exportSubItems = []SubItem{ExportCurrentFile, ExportBatchProcess}

var exportSubLabels = map[SubItem]string{
	ExportCurrentFile:  "Current File",
	ExportBatchProcess: "Batch Process",
}

// Result communicates what the user selected back to the parent.
type Result struct {
	Action    Item
	ThemeName string
	ExportSub SubItem
}

type Model struct {
	Open              bool
	ShouldClose       bool // set when user navigates up from first item
	ShouldNavRight    bool // set when user navigates right with no submenu
	focus             Item
	subOpen           Item // which submenu is open (-1 = none)
	subFocus          int
	SelectedResult    *Result
	width             int
}

func New(w int) Model {
	return Model{
		focus:   LoadSettings,
		subOpen: -1,
		width:   w,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.Open {
		return m, nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.subOpen >= 0 {
		return m.handleSubMenuKey(keyMsg)
	}
	return m.handleMainMenuKey(keyMsg)
}

func (m Model) handleMainMenuKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Up):
		if m.focus == menuItems[0] {
			// At top of menu — close and return focus to File button
			m.Open = false
			m.subOpen = -1
			m.ShouldClose = true
			return m, nil
		}
		for i, item := range menuItems {
			if item == m.focus && i > 0 {
				m.focus = menuItems[i-1]
				return m, nil
			}
		}
	case key.Matches(msg, event.KeyMap.Down):
		for i, item := range menuItems {
			if item == m.focus && i+1 < len(menuItems) {
				m.focus = menuItems[i+1]
				return m, nil
			}
		}
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleSelect()
	case key.Matches(msg, event.KeyMap.Right):
		if hasSubmenu[m.focus] {
			return m.handleSelect()
		}
		// No submenu — navigate right to next button
		m.Open = false
		m.subOpen = -1
		m.ShouldNavRight = true
	case key.Matches(msg, event.KeyMap.Tab):
		for i, item := range menuItems {
			if item == m.focus && i+1 < len(menuItems) {
				m.focus = menuItems[i+1]
				return m, nil
			}
		}
	case key.Matches(msg, event.KeyMap.Esc), key.Matches(msg, event.KeyMap.Left):
		m.Open = false
		m.subOpen = -1
		m.ShouldClose = true
	}
	return m, nil
}

func (m Model) handleSelect() (Model, tea.Cmd) {
	switch m.focus {
	case Theme:
		m.subOpen = Theme
		m.subFocus = 0
	case Export:
		m.subOpen = Export
		m.subFocus = 0
	case LoadSettings, SaveSettings, Quit:
		m.SelectedResult = &Result{Action: m.focus}
		m.Open = false
		m.subOpen = -1
	}
	return m, nil
}

func (m Model) handleSubMenuKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	items := m.currentSubItems()

	switch {
	case key.Matches(msg, event.KeyMap.Up):
		if m.subFocus > 0 {
			m.subFocus--
		}
	case key.Matches(msg, event.KeyMap.Down):
		if m.subFocus+1 < len(items) {
			m.subFocus++
		}
	case key.Matches(msg, event.KeyMap.Tab):
		if m.subFocus+1 < len(items) {
			m.subFocus++
		}
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleSubSelect(items)
	case key.Matches(msg, event.KeyMap.Esc), key.Matches(msg, event.KeyMap.Left):
		m.subOpen = -1
	}
	return m, nil
}

func (m Model) handleSubSelect(items []SubItem) (Model, tea.Cmd) {
	if m.subFocus >= len(items) {
		return m, nil
	}
	selected := items[m.subFocus]

	switch m.subOpen {
	case Theme:
		m.SelectedResult = &Result{
			Action:    Theme,
			ThemeName: themeSubToName[selected],
		}
	case Export:
		m.SelectedResult = &Result{
			Action:    Export,
			ExportSub: selected,
		}
	}
	m.Open = false
	m.subOpen = -1
	return m, nil
}

func (m Model) currentSubItems() []SubItem {
	switch m.subOpen {
	case Theme:
		return themeSubItems
	case Export:
		return exportSubItems
	}
	return nil
}

func (m *Model) Toggle() {
	m.Open = !m.Open
	if m.Open {
		m.focus = LoadSettings
		m.subOpen = -1
		m.SelectedResult = nil
		m.ShouldClose = false
		m.ShouldNavRight = false
	}
}

func (m *Model) Close() {
	m.Open = false
	m.subOpen = -1
}

// MainMenuView returns just the main dropdown rendering.
func (m Model) MainMenuView() string {
	if !m.Open {
		return ""
	}
	return m.renderMainMenu()
}

// SubMenuView returns the submenu rendering, or empty if none is open.
func (m Model) SubMenuView() string {
	if !m.Open || m.subOpen < 0 {
		return ""
	}
	return m.renderSubMenu()
}

// MainMenuWidth returns the rendered width of the main menu.
func (m Model) MainMenuWidth() int {
	return lipgloss.Width(m.renderMainMenu())
}

// SubMenuYOffset returns how many rows down from the dropdown top
// the submenu should shift to align with its parent item.
func (m Model) SubMenuYOffset() int {
	for i, item := range menuItems {
		if Item(m.subOpen) == item {
			return i + 1 // +1 for the top border
		}
	}
	return 0
}

func (m Model) renderMainMenu() string {
	menuWidth := 20
	// Use heavy border when main menu is focused (no submenu open)
	focused := m.subOpen < 0

	lines := make([]string, len(menuItems))
	for i, item := range menuItems {
		itemStyle := style.BgStyle().Width(menuWidth).Foreground(style.DimmedColor1).PaddingLeft(1)
		if item == m.focus {
			itemStyle = style.BgStyle().Width(menuWidth).Foreground(style.SelectedColor1).PaddingLeft(1)
		}
		lines[i] = itemStyle.Render(menuLabels[item])
	}
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	borderType := lipgloss.RoundedBorder()
	if focused {
		borderType = style.HeavyBorder()
	}
	border := style.BgStyle().
		Border(borderType).
		BorderForeground(style.NormalColor1).
		BorderBackground(style.ActiveTheme.Bg)
	return border.Render(content)
}

func (m Model) renderSubMenu() string {
	var items []SubItem
	var labels map[SubItem]string

	switch m.subOpen {
	case Theme:
		items = themeSubItems
		labels = themeSubLabels
	case Export:
		items = exportSubItems
		labels = exportSubLabels
	default:
		return ""
	}

	menuWidth := 26
	lines := make([]string, len(items))
	for i, item := range items {
		itemStyle := style.BgStyle().Width(menuWidth).Foreground(style.DimmedColor1).PaddingLeft(1)
		if i == m.subFocus {
			itemStyle = style.BgStyle().Width(menuWidth).Foreground(style.SelectedColor1).PaddingLeft(1)
		}
		lines[i] = itemStyle.Render(labels[item])
	}
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Submenu always has heavy border when open (it's focused)
	border := style.BgStyle().
		Border(style.HeavyBorder()).
		BorderForeground(style.NormalColor1).
		BorderBackground(style.ActiveTheme.Bg)
	return border.Render(content)
}
