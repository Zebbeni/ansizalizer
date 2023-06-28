package lospec

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
	"github.com/Zebbeni/ansizalizer/style"
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

	width             int
	didInitializeList bool
}

func New(w int) Model {
	return Model{
		focus: CountForm,

		countInput: newInput(CountForm, "8"),
		tagInput:   newInput(TagForm, "gameboy"),
		filterType: FilterMin,
		sortType:   SortDownloads,

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
	sortButtons := m.drawSortButtons()

	titleString := fmt.Sprintf("%d results found\npage %d of %d", len(m.paletteList.Items()), m.paletteList.Paginator.Page, m.paletteList.Paginator.TotalPages)
	title := style.DimmedTitle.Copy().Width(m.width).Height(2).AlignHorizontal(lipgloss.Center).Padding(1, 0, 1, 0).Render(titleString)
	paletteList := m.paletteList.View()
	if len(m.paletteList.Items()) == 0 {
		paletteList = ""
	}
	return lipgloss.JoinVertical(lipgloss.Top, colorFilters, tagInput, sortButtons, title, paletteList)
}

func (m Model) LoadInitial() (Model, tea.Cmd) {
	return m.searchLospec(0)
}

func (m Model) GetCurrent() palette.Model {
	return m.palette
}
