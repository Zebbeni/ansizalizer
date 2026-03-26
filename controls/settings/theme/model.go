package theme

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type State int

const (
	LightOnTransparent State = iota
	LightOnDark
	DarkOnLight
	DarkOnTransparent
	LightOnDarkPaletted
	DarkOnLightPaletted
)

var themeOrder = []State{LightOnTransparent, LightOnDark, LightOnDarkPaletted, DarkOnTransparent, DarkOnLight, DarkOnLightPaletted}

var stateToTheme = map[State]style.ThemeName{
	LightOnTransparent:  style.LightOnTransparent,
	LightOnDark:         style.LightOnDark,
	DarkOnLight:         style.DarkOnLight,
	DarkOnTransparent:   style.DarkOnTransparent,
	LightOnDarkPaletted: style.LightOnDarkPaletted,
	DarkOnLightPaletted: style.DarkOnLightPaletted,
}

var themeToState = map[style.ThemeName]State{
	style.LightOnTransparent:  LightOnTransparent,
	style.LightOnDark:         LightOnDark,
	style.DarkOnLight:         DarkOnLight,
	style.DarkOnTransparent:   DarkOnTransparent,
	style.LightOnDarkPaletted: LightOnDarkPaletted,
	style.DarkOnLightPaletted: DarkOnLightPaletted,
}

type Model struct {
	focus  State
	active State

	ShouldClose bool
	IsActive    bool
	width       int
}

func New(w int) Model {
	return Model{
		focus:       LightOnTransparent,
		active:      LightOnTransparent,
		ShouldClose: false,
		IsActive:    false,
		width:       w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		}
	}
	return m, nil
}

func (m Model) View() string {
	buttons := make([]string, len(themeOrder))
	for i, t := range themeOrder {
		name := style.ThemeNames[stateToTheme[t]]
		btnStyle := style.NormalButtonNode
		if m.IsActive && m.focus == t {
			btnStyle = style.FocusButtonNode
		} else if m.active == t {
			btnStyle = style.ActiveButtonNode
		}
		buttons[i] = btnStyle.Render(name)
	}
	return lipgloss.JoinVertical(lipgloss.Left, buttons...)
}

func (m Model) Summary() string {
	return style.ThemeNames[stateToTheme[m.active]]
}

func (m *Model) ResetFocus() {
	m.focus = m.active
}

func (m Model) ThemeName() string {
	return style.ThemeNames[stateToTheme[m.active]]
}

func (m *Model) SetThemeByName(name string) {
	if themeName, ok := style.ThemeFromName[name]; ok {
		if state, ok := themeToState[themeName]; ok {
			m.active = state
			m.focus = state
			style.SetTheme(themeName)
		}
	}
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	style.SetTheme(stateToTheme[m.active])
	return m, event.StartRenderToViewCmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Down):
		for i, t := range themeOrder {
			if t == m.focus && i+1 < len(themeOrder) {
				m.focus = themeOrder[i+1]
				return m, nil
			}
		}
		m.ShouldClose = true
	case key.Matches(msg, event.KeyMap.Up):
		for i, t := range themeOrder {
			if t == m.focus {
				if i == 0 {
					m.ShouldClose = true
				} else {
					m.focus = themeOrder[i-1]
				}
				return m, nil
			}
		}
	}
	return m, nil
}
