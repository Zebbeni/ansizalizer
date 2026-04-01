package dithering

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
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
	strengthInput    numberinput.Model
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

func newStrengthInput() numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:    "Strength: ",
		CharLimit: 5,
		IsFloat:   true,
		Min:       numberinput.FloatPtr(0),
		Default:   1.0,
	})
	promptSt := lipgloss.NewStyle().Width(12)
	styles := m.Styles()
	styles.Focused.Prompt = promptSt
	styles.Blurred.Prompt = promptSt
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	m.SetStyles(styles)
	return m
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
	pad := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
	parts := []string{pad.Render(m.drawDitheringOptions())}
	if m.doDithering {
		parts = append(parts, pad.Render(m.drawSerpentineOptions()))
		parts = append(parts, pad.Render(m.drawStrength()))
		parts = append(parts, m.drawModeTabs())
	}
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return content
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
	return m.strengthInput.Float32Value()
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

func (m *Model) selectMatrixListItem(mt MatrixType) {
	for i, item := range m.matrixList.Items() {
		if mi, ok := item.(matrixItem); ok && mi.Type == mt {
			m.matrixList.Select(i)
			return
		}
	}
}

func (m *Model) selectClusteredDotListItem(ct ClusteredDotType) {
	for i, item := range m.clusteredDotList.Items() {
		if ci, ok := item.(clusteredDotItem); ok && ci.Type == ct {
			m.clusteredDotList.Select(i)
			return
		}
	}
}

func (m *Model) resetMatrixListCursor() {
	for mt, matrix := range errorDiffMatrixMap {
		if matrixEqual(matrix, m.matrix) {
			m.selectMatrixListItem(mt)
			return
		}
	}
}

func (m *Model) resetClusteredDotListCursor() {
	for ct, cdm := range clusteredDotMatrixMap {
		if orderedMatrixEqual(cdm, m.clusteredDot) {
			m.selectClusteredDotListItem(ct)
			return
		}
	}
}

func (m *Model) SetConfig(doDithering, doSerpentine bool, matrixName string) {
	m.doDithering = doDithering
	m.doSerpentine = doSerpentine
	for mt, name := range matrixNameMap {
		if name == matrixName {
			if matrix, ok := errorDiffMatrixMap[mt]; ok {
				m.matrix = matrix
				m.selectMatrixListItem(mt)
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
				m.selectMatrixListItem(mt)
				break
			}
		}
	}
	for ct, name := range clusteredDotNameMap {
		if name == clusteredDotName {
			if cdm, ok := clusteredDotMatrixMap[ct]; ok {
				m.clusteredDot = cdm
				m.selectClusteredDotListItem(ct)
				break
			}
		}
	}
	if bayerSize > 0 {
		m.bayerSize = uint(bayerSize)
	}
	if strength > 0 {
		m.strengthInput.SetFloatValue(float64(strength), -1)
	}
}
