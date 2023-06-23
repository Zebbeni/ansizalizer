package lospec

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/textinput"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
)

type State int

const (
	CountForm State = iota
	TagForm
	FilterExact
	FilterMax
	FilterMin
	SortAlphabetical
	SortDownloads
	SortNewest
	List
)

type Model struct {
	focus  State
	active State

	countInput textinput.Model
	tagInput   textinput.Model
	filterType State
	sortType   State

	paletteList            list.Model
	palettes               []list.Item
	palette                palette.Model
	isPaletteListAllocated bool
	highestPageRequested   int
	requestID              int

	ShouldClose   bool
	ShouldUnfocus bool
	IsActive      bool

	width int
}

func New(w int) Model {
	return Model{
		focus: CountForm,

		countInput: newInput(CountForm, "32"),
		tagInput:   newInput(TagForm, ""),
		filterType: FilterExact,
		sortType:   SortAlphabetical,

		isPaletteListAllocated: false,
		highestPageRequested:   0,
		requestID:              0,

		ShouldClose:   false,
		ShouldUnfocus: false,
		IsActive:      false,

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case CountForm:
		if m.countInput.Focused() {
			return m.handleCountFormUpdate(msg)
		}
	case TagForm:
		if m.tagInput.Focused() {
			return m.handleTagFormUpdate(msg)
		}
	}

	switch msg := msg.(type) {
	case event.LospecResponseMsg:
		return m.handleLospecResponse(msg)
	case tea.KeyMsg:
		if m.focus == List {
			return m.handleListUpdate(msg)
		}
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

// View draws a control panel like this:
//
// Colors ___ |Exact Max Min
// Tag _____________________
// Sort By |A-Z Downloads New
//
// (Palette List)
// <palette name>
// <preview>
// <...>
// <...>
// ..
func (m Model) View() string {
	colorsInput := m.drawColorsInput()
	filters := m.drawFilterButtons()
	colorFilters := lipgloss.JoinHorizontal(lipgloss.Left, colorsInput, filters)
	tagInput := m.drawTagInput()
	paletteList := m.drawPaletteList()

	return lipgloss.JoinVertical(lipgloss.Top, colorFilters, tagInput, paletteList)
}

func (m Model) GetCurrent() palette.Model {
	return m.palette
}
