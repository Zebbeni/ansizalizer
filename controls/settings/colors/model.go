package colors

import (
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/controls/settings/colors/adaptive"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors/limited"
)

type State int

// None consists of a few different components that are shown or hidden
// depending on which toggles have been set on / off. The Model state indicates
// which component is currently focused. From top to bottom the components are:

// 1) Limited (on/off)
// 2) Load (Name) (if Limited) -> [Enter] displays Load menu
// 3) Dithering (on/off) (if Limited)
// 4) Serpentine (on/off) (if Dithering)
// 5) Matrix (Name) (if Dithering) -> [Enter] displays to Matrix menu

// These can all be part of a single list, but we need to onSelect the list items

const (
	Full State = iota
	Create
	Load
	Lospec
	CreateControls
	LoadControls
	LospecControls
)

type Model struct {
	selected State
	focus    State // the component taking input
	controls State

	Adaptive adaptive.Model
	Palette  limited.Model

	ShouldClose      bool
	ShouldDeactivate bool

	IsActive bool

	width int
}

func New(w int) Model {
	m := Model{
		selected:         Full,
		focus:            Full,
		controls:         Full,
		Adaptive:         adaptive.New(w),
		Palette:          limited.New(w),
		ShouldClose:      false,
		ShouldDeactivate: false,
		IsActive:         false,
		width:            w,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.focus {
	case CreateControls:
		return m.handleAdaptiveUpdate(msg)
	case LoadControls:
		return m.handlePaletteUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

func (m Model) View() string {
	buttons := m.drawButtons()
	if m.IsActive == false {
		return buttons
	}

	var controls string
	switch m.controls {
	case Create:
		controls = m.Adaptive.View()
	case Load:
		controls = m.Palette.View()
	}
	if len(controls) == 0 {
		return buttons
	}

	return lipgloss.JoinVertical(lipgloss.Top, buttons, controls)
}

func (m Model) IsLimited() bool {
	return m.selected != Full
}

func (m Model) IsDithered() bool {
	return false
}

func (m Model) IsSerpentine() bool {
	return true
}

func (m Model) Matrix() dither.ErrorDiffusionMatrix {
	return dither.FloydSteinberg
}

func (m Model) IsAdaptive() bool {
	return m.selected == Create
}

func (m Model) IsPaletted() bool {
	return m.selected == Load
}

func (m Model) GetCurrentPalette() color.Palette {
	if m.selected == Load {
		return m.Palette.GetCurrent()
	}
	return m.Adaptive.Palette
}
