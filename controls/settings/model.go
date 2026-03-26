package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings/adjust"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced"
	"github.com/Zebbeni/ansizalizer/controls/settings/alpha"
	"github.com/Zebbeni/ansizalizer/controls/settings/animation"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors"
	"github.com/Zebbeni/ansizalizer/controls/settings/saveload"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
	"github.com/Zebbeni/ansizalizer/controls/settings/theme"
	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	active   State
	focus    State
	expanded map[State]bool

	Colors     colors.Model
	Characters characters.Model
	Size       size.Model
	Adjust     adjust.Model
	Advanced   advanced.Model
	Animation  animation.Model
	Alpha      alpha.Model
	Theme      theme.Model
	SaveLoad   saveload.Model

	ShouldUnfocus bool
	ShouldClose   bool

	width int
}

func New(w int, dirs prefs.Dirs) Model {
	return Model{
		active:   None,
		focus:    SaveLoad,
		expanded: make(map[State]bool),

		Colors:     colors.New(w, dirs.PaletteLoad),
		Characters: characters.New(w - 2),
		Size:       size.New(),
		Adjust:     adjust.New(w - 2),
		Advanced:   advanced.New(w - 2),
		Animation:  animation.New(),
		Alpha:      alpha.New(w - 2),
		Theme:      theme.New(w - 2),
		SaveLoad:   saveload.New(w-2, dirs.SaveDir, dirs.LoadFile),

		ShouldUnfocus: false,
		ShouldClose:   false,

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) DebugState() string {
	expandedList := ""
	for s := range m.expanded {
		expandedList += " " + stateTitles[s]
	}
	return fmt.Sprintf("Settings: active=%s focus=%s expanded=[%s]\n  Theme: %s\n  Colors: trueColor=%v palette=%s adapt=%v\n  Animation: visible=%v\n  Alpha: visible=%v",
		stateTitles[m.active], stateTitles[m.focus], expandedList,
		m.Theme.ThemeName(),
		!m.Colors.IsLimited(), m.Colors.GetCurrentPalette().Name(), m.Colors.AdaptToPalette(),
		m.Animation.AnimatedImage,
		m.Alpha.AlphaImage,
	)
}

func (m *Model) syncPalette() {
	m.Advanced.ShowDithering = m.Colors.IsLimited()
	if m.Colors.IsLimited() {
		style.PaletteColors = m.Colors.GetCurrentPalette().Colors()
	} else {
		style.PaletteColors = nil
	}
	if style.ActiveTheme.Name == style.LightOnDarkPaletted || style.ActiveTheme.Name == style.DarkOnLightPaletted {
		style.SetTheme(style.ActiveTheme.Name)
		m.Colors.RefreshStyles()
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.Advanced.ShowDithering = m.Colors.IsLimited()

	switch m.active {
	case Colors:
		m, cmd := m.handleColorsUpdate(msg)
		m.syncPalette()
		return m, cmd
	case Characters:
		return m.handleCharactersUpdate(msg)
	case Size:
		return m.handleSizeUpdate(msg)
	case Adjust:
		return m.handleAdjustUpdate(msg)
	case Advanced:
		return m.handleAdvancedUpdate(msg)
	case Animation:
		return m.handleAnimationUpdate(msg)
	case Alpha:
		return m.handleAlphaUpdate(msg)
	case Theme:
		return m.handleThemeUpdate(msg)
	case SaveLoad:
		return m.handleSaveLoadUpdate(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKeyMsg(keyMsg)
}

func (m Model) View() string {
	col := m.renderCollapsible(m.Colors.View(), m.Colors.Summary(), Colors)
	char := m.renderCollapsible(m.Characters.View(), m.Characters.Summary(), Characters)
	siz := m.renderCollapsible(m.Size.View(), m.Size.Summary(), Size)
	adj := m.renderCollapsible(m.Adjust.View(), m.Adjust.Summary(), Adjust)
	sam := m.renderCollapsible(m.Advanced.View(), m.Advanced.Summary(), Advanced)
	thm := m.renderCollapsible(m.Theme.View(), m.Theme.Summary(), Theme)
	sl := m.renderCollapsible(m.SaveLoad.View(), m.SaveLoad.Summary(), SaveLoad)
	parts := []string{sl, thm, col, char, siz, adj, sam}
	if m.Animation.AnimatedImage {
		parts = append(parts, m.renderCollapsible(m.Animation.View(), m.Animation.Summary(), Animation))
	}
	if m.Alpha.AlphaImage {
		parts = append(parts, m.renderCollapsible(m.Alpha.View(), m.Alpha.Summary(), Alpha))
	}

	return lipgloss.JoinVertical(lipgloss.Top, parts...)
}

func (m Model) renderCollapsible(content, summary string, state State) string {
	if m.active == state || m.expanded[state] {
		return m.renderWithBorder(content, state, false)
	}
	dimSummary := style.BgStyle().Foreground(style.ExtraDimColor).Italic(true).Width(m.width - 2).AlignHorizontal(lipgloss.Center).Padding(0, 1).Render(summary)
	return m.renderWithBorder(dimSummary, state, true)
}
