package palettes

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/settings/palettes/adaptive"
	"github.com/Zebbeni/ansizalizer/controls/settings/palettes/loader"
	"github.com/Zebbeni/ansizalizer/controls/settings/palettes/lospec"
	"github.com/Zebbeni/ansizalizer/palette"
)

type State int

// None consists of a few different components that are shown or hidden
// depending on which toggles have been set on / off. The Model state indicates
// which component is currently focused. From top to bottom the components are:

// 1) Limited (on/off)
// 2) Loader (Name) (if Limited) -> [Enter] displays Loader menu
// 3) Dithering (on/off) (if Limited)
// 4) Serpentine (on/off) (if Dithering)
// 5) Matrix (Name) (if Dithering) -> [Enter] displays to Matrix menu

// These can all be part of a single list, but we need to onSelect the list items

const (
	Adapt State = iota
	Load
	Lospec
	AdaptiveControls
	LoadControls
	LospecControls
)

type Model struct {
	selected State
	focus    State // the component taking input
	controls State

	Adapter adaptive.Model
	Loader  loader.Model
	Lospec  lospec.Model

	ShouldClose bool

	IsActive bool

	width int
}

func New(w int, paletteDir string) Model {
	innerW := w - 2
	loaderModel := loader.New(innerW)
	if paletteDir != "" {
		loaderModel = loader.NewAtDir(paletteDir, innerW)
	}
	m := Model{
		selected:    Lospec,
		focus:       Lospec,
		controls:    Lospec,
		Adapter:     adaptive.New(innerW),
		Loader:      loaderModel,
		Lospec:      lospec.New(innerW),
		ShouldClose: false,
		IsActive:    false,
		width:       w,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.focus {
	case AdaptiveControls:
		return m.handleAdaptiveUpdate(msg)
	case LoadControls:
		return m.handleLoaderUpdate(msg)
	case LospecControls:
		return m.handleLospecUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

func (m Model) View() string {
	return m.drawTabs()
}

func (m *Model) SetSelected(state State) {
	m.selected = state
	m.controls = state
	m.focus = state
}

func (m Model) IsLospecFocused() bool {
	return m.controls == Lospec
}

func (m Model) IsAdaptive() bool {
	return m.selected == Adapt
}

func (m Model) IsPaletted() bool {
	return m.selected == Load
}

func (m *Model) RefreshStyles() {
	m.Lospec.RefreshStyles()
}

func (m Model) GetCurrentPalette() palette.Model {
	switch m.selected {
	case Load:
		return m.Loader.GetCurrent()
	case Adapt:
		return m.Adapter.GetCurrent()
	case Lospec:
		return m.Lospec.GetCurrent()
	}
	return palette.Model{}
}
