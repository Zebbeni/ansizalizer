package export

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/export/destination"
	"github.com/Zebbeni/ansizalizer/controls/export/source"
)

type State int

const (
	None State = iota
	Source
	Destination
	Format
)

var (
	stateTitles = map[State]string{
		Source:      "Source",
		Destination: "Destination",
		Format:      "Format",
	}
)

type Model struct {
	active State
	focus  State

	Source      source.Model
	Destination destination.Model

	ShouldClose   bool
	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	return Model{
		focus:         Source,
		active:        None,
		Source:        source.New(w - 2),
		Destination:   destination.New(w - 2),
		ShouldClose:   false,
		ShouldUnfocus: false,
		width:         w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case Source:
		return m.handleSourceUpdate(msg)
	case Destination:
		return m.handleDestinationUpdate(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKeyMsg(keyMsg)
}

func (m Model) View() string {
	src := m.renderWithBorder(m.Source.View(), Source)
	dst := m.renderWithBorder(m.Destination.View(), Destination)
	return lipgloss.JoinVertical(lipgloss.Left, src, dst)
}
