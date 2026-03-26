package viewer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type Model struct {
	imgString   string
	wrappedView string
	settings    settings.Model

	// Animation state
	frames       []string
	delay        time.Duration
	currentFrame int
	isAnimating  bool
	generation   int

	// Panel dimensions for box wrapping
	panelW, panelH int

	WaitingOnRender bool
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case event.FinishRenderToViewMsg:
		return m.handleFinishRenderMsg(msg)
	case event.FinishRenderGIFToViewMsg:
		return m.handleFinishRenderGIFMsg(msg)
	case event.AnimationTickMsg:
		return m.handleAnimationTick(msg)
	}
	return m, nil
}

func (m *Model) SetSize(w, h int) {
	if m.panelW != w || m.panelH != h {
		m.panelW = w
		m.panelH = h
		m.wrappedView = m.buildWrappedView()
	}
}

func (m Model) buildWrappedView() string {
	imgString := m.imgString
	if imgString == "" {
		return ""
	}

	imgWidth, imgHeight := lipgloss.Size(imgString)
	imgViewer := imgString

	if imgHeight > 1 && imgWidth > 4 {
		boxLabelRenderer := style.BoxWithLabel{
			BoxStyle:   style.BgStyle().BorderForeground(style.ExtraDimColor).Border(lipgloss.RoundedBorder()).BorderBackground(style.ActiveTheme.Bg),
			LabelStyle: style.BgStyle().Foreground(style.ExtraDimColor).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Bottom),
		}
		imgViewer = boxLabelRenderer.Render(fmt.Sprintf("%dx%d", imgWidth, imgHeight), imgString, imgWidth)
	}

	displayHeight := 3
	renderVP := viewport.New(m.panelW-2, m.panelH-displayHeight-2)

	vpRightStyle := style.BgStyle().Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	rightContent := vpRightStyle.Copy().Width(m.panelW - 2).Height(m.panelH - 4).Render(imgViewer)
	renderVP.SetContent(rightContent)
	renderVP.Style = style.BgStyle()

	content := renderVP.View()
	return style.NormalButton.Copy().BorderForeground(style.DimmedColor1).BorderBackground(style.ActiveTheme.Bg).Render(content)
}

func (m Model) IsAnimating() bool {
	return m.isAnimating && len(m.frames) > 1
}

func (m Model) Frames() []string {
	return m.frames
}

func (m Model) View() string {
	return m.wrappedView
}

func (m Model) ImgString() string {
	return m.imgString
}
