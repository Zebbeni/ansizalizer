package dithering

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	DitherOn State = iota
	DitherOff
	SerpentineOn
	SerpentineOff
	StrengthForm
	ModeMatrix
	ModeBayer
	ModeClusteredDot
	MatrixList
	BayerSize2
	BayerSize4
	BayerSize8
	BayerSize16
	ClusteredDotList
)

type DitherMode int

const (
	Matrix DitherMode = iota
	Bayer
	ClusteredDot
)

type Model struct {
	focus        State
	active       State
	modeControls State // which tab content to show

	doDithering  bool
	doSerpentine bool
	ditherMode   DitherMode

	matrix           dither.ErrorDiffusionMatrix
	matrixList       list.Model
	bayerSize        uint
	strengthInput    textinput.Model
	clusteredDot     dither.OrderedDitherMatrix
	clusteredDotList list.Model

	IsActive    bool
	ShouldClose bool

	width int
}

func New(w int) Model {
	return Model{
		focus:            DitherOff,
		active:           DitherOff,
		modeControls:     ModeMatrix,
		doDithering:      false,
		doSerpentine:     false,
		ditherMode:       Matrix,
		matrix:           dither.FloydSteinberg,
		matrixList:       newMatrixMenu(w),
		bayerSize:        4,
		strengthInput:    newStrengthInput(),
		clusteredDot:     dither.ClusteredDot4x4,
		clusteredDotList: newClusteredDotMenu(w),
		ShouldClose:      false,
		IsActive:         false,
		width:            w,
	}
}

func newStrengthInput() textinput.Model {
	input := textinput.New()
	input.Prompt = "Strength: "
	input.PromptStyle = lipgloss.NewStyle().Width(12)
	input.CharLimit = 5
	input.SetValue("1.0")
	return input
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Handle focused text input
	if m.active == StrengthForm && m.strengthInput.Focused() {
		return m.handleStrengthUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle list navigation when focused on a list
		if m.focus == MatrixList {
			return m.handleMatrixListUpdate(msg)
		}
		if m.focus == ClusteredDotList {
			return m.handleClusteredDotListUpdate(msg)
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

func (m Model) View() string {
	parts := []string{m.drawDitheringOptions()}
	if m.doDithering {
		parts = append(parts, m.drawSerpentineOptions())
		parts = append(parts, m.drawStrength())
		parts = append(parts, m.drawModeTabs())
	}
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return lipgloss.NewStyle().Padding(0, 1).Render(content)
}

func (m Model) Settings() (bool, bool, dither.ErrorDiffusionMatrix) {
	return m.doDithering, m.doSerpentine, m.matrix
}

func (m Model) DitherModeValue() DitherMode {
	return m.ditherMode
}

func (m Model) BayerSize() uint {
	if m.bayerSize == 0 {
		return 4
	}
	return m.bayerSize
}

func (m Model) Strength() float32 {
	val, err := strconv.ParseFloat(m.strengthInput.Value(), 32)
	if err != nil || val <= 0 {
		return 1.0
	}
	return float32(val)
}

func (m Model) ClusteredDotMatrix() dither.OrderedDitherMatrix {
	return m.clusteredDot
}

func (m Model) MatrixName() string {
	for mt, matrix := range errorDiffMatrixMap {
		if matrixEqual(matrix, m.matrix) {
			return matrixNameMap[mt]
		}
	}
	return "FloydSteinberg"
}

func (m Model) ClusteredDotName() string {
	for ct, cdm := range clusteredDotMatrixMap {
		if orderedMatrixEqual(cdm, m.clusteredDot) {
			return clusteredDotNameMap[ct]
		}
	}
	return "ClusteredDot4x4"
}

func matrixEqual(a, b dither.ErrorDiffusionMatrix) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func orderedMatrixEqual(a, b dither.OrderedDitherMatrix) bool {
	if len(a.Matrix) != len(b.Matrix) || a.Max != b.Max {
		return false
	}
	for i := range a.Matrix {
		if len(a.Matrix[i]) != len(b.Matrix[i]) {
			return false
		}
		for j := range a.Matrix[i] {
			if a.Matrix[i][j] != b.Matrix[i][j] {
				return false
			}
		}
	}
	return true
}

func (m *Model) SetConfig(doDithering, doSerpentine bool, matrixName string) {
	m.doDithering = doDithering
	m.doSerpentine = doSerpentine
	for mt, name := range matrixNameMap {
		if name == matrixName {
			if matrix, ok := errorDiffMatrixMap[mt]; ok {
				m.matrix = matrix
				return
			}
		}
	}
}

func (m *Model) SetFullConfig(doDithering, doSerpentine bool, mode DitherMode, matrixName, clusteredDotName string, bayerSize int, strength float32) {
	m.doDithering = doDithering
	m.doSerpentine = doSerpentine
	m.ditherMode = mode
	switch mode {
	case Matrix:
		m.modeControls = ModeMatrix
	case Bayer:
		m.modeControls = ModeBayer
	case ClusteredDot:
		m.modeControls = ModeClusteredDot
	}
	for mt, name := range matrixNameMap {
		if name == matrixName {
			if matrix, ok := errorDiffMatrixMap[mt]; ok {
				m.matrix = matrix
				break
			}
		}
	}
	for ct, name := range clusteredDotNameMap {
		if name == clusteredDotName {
			if cdm, ok := clusteredDotMatrixMap[ct]; ok {
				m.clusteredDot = cdm
				break
			}
		}
	}
	if bayerSize > 0 {
		m.bayerSize = uint(bayerSize)
	}
	if strength > 0 {
		m.strengthInput.SetValue(strconv.FormatFloat(float64(strength), 'f', -1, 32))
	}
}
