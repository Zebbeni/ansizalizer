package colors

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings/palettes"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
)

type State int

const (
	UsePalette State = iota
	UseTrueColor
	Palette
)

type Model struct {
	focus           State
	mode            State
	width           int
	PaletteControls palettes.Model

	IsActive    bool
	ShouldClose bool
}

func New(w int) Model {
	return Model{
		focus:           UseTrueColor,
		mode:            UseTrueColor,
		width:           w,
		PaletteControls: palettes.New(w),
		IsActive:        false,
		ShouldClose:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.focus {
	case Palette:
		return m.handlePaletteUpdate(msg)
	}

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
	paletteToggles := m.drawPaletteToggles()
	if m.mode == UseTrueColor {
		return paletteToggles
	}

	paletteTabs := m.PaletteControls.View()
	return lipgloss.JoinVertical(lipgloss.Left, paletteToggles, paletteTabs)
}

// GetSelected returns isPaletted, isAdaptive, and the palette (if applicable)
func (m Model) GetSelected() (bool, bool, palette.Model) {
	colorPalette := m.PaletteControls.GetCurrentPalette()

	if m.mode == UseTrueColor {
		return true, false, colorPalette
	}

	return false, m.PaletteControls.IsAdaptive(), colorPalette
}

func (m Model) GetCurrentPalette() palette.Model {
	return m.PaletteControls.GetCurrentPalette()
}

func (m Model) IsLimited() bool {
	return m.mode == UsePalette
}
