package colors

import (
	"fmt"
	"image/color"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings/palettes"
	"github.com/Zebbeni/ansizalizer/controls/settings/palettes/loader"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
)

type State int

const (
	UsePalette State = iota
	UseTrueColor
	AdaptOn
	AdaptOff
	Palette
)

type Model struct {
	focus           State
	mode            State
	adaptToPalette  bool
	width           int
	PaletteControls palettes.Model

	IsActive    bool
	ShouldClose bool
}

func New(w int, paletteDir string) Model {
	return Model{
		focus:           UseTrueColor,
		mode:            UseTrueColor,
		width:           w,
		PaletteControls: palettes.New(w-2, paletteDir),
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
	adaptToggle := m.drawAdaptToggle()
	if m.mode == UseTrueColor {
		return lipgloss.JoinVertical(lipgloss.Left, paletteToggles, adaptToggle)
	}

	paletteTabs := m.PaletteControls.View()
	return lipgloss.JoinVertical(lipgloss.Left, paletteToggles, adaptToggle, paletteTabs)
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

// NewPaletteMode returns a Model configured to use the given palette.
func NewPaletteMode(p color.Palette, w int) Model {
	m := New(w, "")
	m.mode = UsePalette
	m.PaletteControls.Loader = loader.NewWithPalette(p, w)
	return m
}

func (m Model) AdaptToPalette() bool {
	return m.adaptToPalette
}

func (m *Model) SetAdaptToPalette(adapt bool) {
	m.adaptToPalette = adapt
}

func (m *Model) SetMode(useTrueColor bool) {
	if useTrueColor {
		m.mode = UseTrueColor
	} else {
		m.mode = UsePalette
	}
}

func (m *Model) SetPalette(name string, colors color.Palette) {
	m.mode = UsePalette
	m.PaletteControls.Loader = loader.NewWithNamedPalette(name, colors, m.width-2)
	m.PaletteControls.SetSelected(palettes.Load)
}

func (m *Model) ResetFocus() {
	m.focus = UseTrueColor
}

func (m Model) Summary() string {
	if m.mode == UseTrueColor {
		return "Mode: True Color"
	}
	p := m.PaletteControls.GetCurrentPalette()
	count := len(p.Colors())
	summary := fmt.Sprintf("Mode: Palette (%d color)", count)
	name := p.Name()
	if name != "" {
		summary += "\nName: " + name
	}
	if m.adaptToPalette {
		summary += "\nAdapt: On"
	}
	return summary
}
