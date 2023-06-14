package settings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
	"github.com/Zebbeni/ansizalizer/controls/settings/colors"
	"github.com/Zebbeni/ansizalizer/controls/settings/sampling"
	"github.com/Zebbeni/ansizalizer/controls/settings/size"
)

type Model struct {
	active State
	focus  State

	Colors     colors.Model
	Characters characters.Model
	Size       size.Model
	Sampling   sampling.Model

	ShouldUnfocus bool
	ShouldClose   bool

	width int
}

func New(w int) Model {
	return Model{
		active: None,
		focus:  Colors,

		Colors:     colors.New(w),
		Characters: characters.New(),
		Size:       size.New(),
		Sampling:   sampling.New(),

		ShouldUnfocus: false,
		ShouldClose:   false,

		width: w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case Colors:
		return m.handleColorsUpdate(msg)
	case Characters:
		return m.handleCharactersUpdate(msg)
	case Size:
		return m.handleSizeUpdate(msg)
	case Sampling:
		return m.handleSamplingUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

func (m Model) View() string {
	col := m.renderWithBorder(m.Colors.View(), Colors)
	char := m.renderWithBorder(m.Characters.View(), Characters)
	siz := m.renderWithBorder(m.Size.View(), Size)
	sam := m.renderWithBorder(m.Sampling.View(), Sampling)

	return lipgloss.JoinVertical(lipgloss.Top, col, char, siz, sam)
}
