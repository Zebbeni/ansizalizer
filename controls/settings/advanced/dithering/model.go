package dithering

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeworld-the-better-one/dither/v2"
)

type State int

const (
	DitherOn State = iota
	DitherOff
	SerpentineOn
	SerpentineOff
	Matrix
)

type Model struct {
	focus State

	doDithering  bool
	doSerpentine bool
	matrix       dither.ErrorDiffusionMatrix

	list list.Model

	IsActive    bool
	ShouldClose bool

	width int
}

func New(w int) Model {
	return Model{
		focus:        DitherOff,
		doDithering:  false,
		doSerpentine: false,
		matrix:       dither.FloydSteinberg,
		list:         newMatrixMenu(w),
		ShouldClose:  false,
		IsActive:     false,
		width:        w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == Matrix {
		return m.handleMatrixListUpdate(msg)
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(keyMsg)
	}
	return m, nil
}

func (m Model) View() string {
	ditheringOpts := m.drawDitheringOptions()
	serpentineOpts := m.drawSerpentineOptions()
	matrixList := m.drawMatrix()
	content := lipgloss.JoinVertical(lipgloss.Left, ditheringOpts, serpentineOpts, matrixList)
	return lipgloss.NewStyle().Padding(0, 1).Render(content)
}

func (m Model) Settings() (bool, bool, dither.ErrorDiffusionMatrix) {
	return m.doDithering, m.doSerpentine, m.matrix
}
