package lospec

import tea "github.com/charmbracelet/bubbletea"

type State int

const (
	CountForm State = iota
	TagForm
	FilterButtons
	SortingButtons
	List
)

type Model struct {
	width int
}

func New(w int) Model {
	return Model{
		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// Colors ___ Tag _____________
// Number Any Exact Max Min
// Sort A-Z Downloads New
//
// <palette name>
// <preview>
// <...>
// <...>
// ..

func (m Model) View() string {
	// draw count input
	// draw tag input
	// join horizontally
	// draw filter buttons
	// draw sorting buttons
	// draw list
	return "Lospec Menu"
}

// https://lospec.com/palette-list/load?colorNumberFilterType=exact&colorNumber=32&tag=&sortingType=alphabetical
