package lospec

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/component/textinput"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	CountForm State = iota
	TagForm
	FilterAny
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
	isPaletteListAllocated bool
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
		filterType: FilterAny,
		sortType:   SortAlphabetical,

		isPaletteListAllocated: false,
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

	if m.focus == List {
		return m.handleListUpdate(msg)
	}

	switch msg := msg.(type) {
	case event.LospecResponseMsg:
		return m.handleLospecResponse(msg)
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

// View draws a control panel like this:
//
// Colors: ___ Tag: _____________
// Number: Any Exact Max Min
// Sort: A-Z Downloads New
//
// (Palette List)
// <palette name>
// <preview>
// <...>
// <...>
// ..
func (m Model) View() string {
	// draw count input
	// draw tag input
	inputs := m.drawInputs()
	// join horizontally
	// draw filter buttons
	filters := m.drawFilterButtons()
	// draw sorting buttons
	// draw list
	paletteList := m.paletteList.View()
	return lipgloss.JoinVertical(lipgloss.Top, inputs, filters, paletteList)
}
