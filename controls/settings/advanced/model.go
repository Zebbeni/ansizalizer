package advanced

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/dithering"
	"github.com/Zebbeni/ansizalizer/controls/settings/advanced/sampling"
	"github.com/Zebbeni/ansizalizer/event"
)

type State int

const (
	Menu State = iota
	Sampling
	Dithering
	SamplingControls
	DitheringControls
)

type Model struct {
	focus         State
	active        State
	activeTab     State
	sampling      sampling.Model
	dithering     dithering.Model
	ShouldClose   bool
	IsActive      bool
	ShowDithering bool
	width         int
}

func New(w int) Model {
	return Model{
		focus:       Sampling,
		active:      Menu,
		activeTab:   Sampling,
		sampling:    sampling.New(w - 2),
		dithering:   dithering.New(w - 2),
		ShouldClose: false,
		IsActive:    false,
		width:       w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.active {
	case SamplingControls:
		return m.handleSamplingUpdate(msg)
	case DitheringControls:
		return m.handleDitheringUpdate(msg)
	}

	switch msg := msg.(type) {
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

func (m Model) View() string {
	return m.drawTabs()
}

func (m Model) SamplingFunction() resize.InterpolationFunction {
	return m.sampling.Function
}

func (m Model) Dithering() (bool, bool, dither.ErrorDiffusionMatrix) {
	return m.dithering.Settings()
}

func (m Model) DitherMode() dithering.DitherMode {
	return m.dithering.DitherModeValue()
}

func (m Model) BayerSize() uint {
	return m.dithering.BayerSize()
}

func (m Model) ClusteredDotMatrix() dither.OrderedDitherMatrix {
	return m.dithering.ClusteredDotMatrix()
}

func (m Model) ClusteredDotName() string {
	return m.dithering.ClusteredDotName()
}

func (m Model) DitherStrength() float32 {
	return m.dithering.Strength()
}

func (m Model) SamplingFunctionName() string {
	return m.sampling.FunctionName()
}

func (m Model) DitherMatrixName() string {
	return m.dithering.MatrixName()
}

func (m *Model) SetConfig(samplingName string, doDither, doSerpentine bool, ditherMode dithering.DitherMode, matrixName, clusteredDotName string, bayerSize int, strength float32) {
	m.sampling.SetFunctionByName(samplingName)
	m.dithering.SetFullConfig(doDither, doSerpentine, ditherMode, matrixName, clusteredDotName, bayerSize, strength)
}

func (m *Model) RefreshStyles() {
	m.sampling.SetListActive(m.sampling.IsActive)
}

func (m *Model) ResetFocus() {
	if m.ShowDithering {
		m.focus = Dithering
		m.activeTab = Dithering
	} else {
		m.focus = Sampling
		m.activeTab = Sampling
	}
	m.active = Menu
}

func (m Model) Summary() string {
	summary := "Sample: " + m.sampling.FunctionName()
	doDither, doSerpentine, _ := m.dithering.Settings()
	if !doDither {
		return summary
	}

	mode := m.dithering.DitherModeValue()
	strength := m.dithering.Strength()
	strengthStr := strconv.FormatFloat(float64(strength), 'f', -1, 32)

	var ditherInfo string
	switch mode {
	case dithering.Matrix:
		ditherInfo = "Dither: " + m.dithering.MatrixName()
	case dithering.Bayer:
		ditherInfo = fmt.Sprintf("Dither: Bayer %d", m.dithering.BayerSize())
	case dithering.ClusteredDot:
		ditherInfo = "Dither: " + m.dithering.ClusteredDotName()
	}
	ditherInfo += " (" + strengthStr + ")"

	if doSerpentine {
		ditherInfo += "\nSerpentine: On"
	}

	return summary + "\n" + ditherInfo
}
