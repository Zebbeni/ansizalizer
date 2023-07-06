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
		focus:  Palette,

		Colors:     colors.New(w),
		Characters: characters.New(w - 2),
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
	case Palette:
		return m.handleColorsUpdate(msg)
	case Characters:
		return m.handleCharactersUpdate(msg)
	case Size:
		return m.handleSizeUpdate(msg)
	case Sampling:
		return m.handleSamplingUpdate(msg)
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKeyMsg(keyMsg)
}

func (m Model) View() string {
	colorCtrls := m.Colors.View()
	charCtrls := m.Characters.View()
	sizeCtrls := m.Size.View()
	sampCtrls := m.Sampling.View()

	col := m.renderWithBorder(colorCtrls, Palette)
	char := m.renderWithBorder(charCtrls, Characters)
	siz := m.renderWithBorder(sizeCtrls, Size)
	sam := m.renderWithBorder(sampCtrls, Sampling)

	return lipgloss.JoinVertical(lipgloss.Top, col, char, siz, sam)
}
