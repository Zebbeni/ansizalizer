package viewer

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

// cropWidth truncates each line of s to maxWidth visual columns,
// handling ANSI escape sequences correctly.
func cropWidth(s string, maxWidth int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if lipgloss.Width(line) > maxWidth {
			lines[i] = ansi.Truncate(line, maxWidth, "")
		}
	}
	return strings.Join(lines, "\n")
}

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
	cropped := false

	// Crop image to fit the viewer width to prevent text wrapping.
	// Account for the border box (2 chars) and the viewport border (2 chars).
	maxWidth := m.panelW - 6
	if imgWidth > maxWidth && maxWidth > 0 {
		imgString = cropWidth(imgString, maxWidth)
		imgWidth = maxWidth
		cropped = true
	}

	// Check if image height exceeds available viewer height
	maxHeight := m.panelH - 7 // display(3) + viewport border(2) + image box border(2)
	if imgHeight > maxHeight && maxHeight > 0 {
		cropped = true
	}

	imgViewer := imgString

	if imgHeight > 1 && imgWidth > 4 {
		label := fmt.Sprintf("%dx%d", imgWidth, imgHeight)
		if cropped {
			label += "(cropped)"
		}
		boxLabelRenderer := style.BoxWithLabel{
			BoxStyle:   style.BgStyle().BorderForeground(style.ExtraDimColor).Border(lipgloss.RoundedBorder()).BorderBackground(style.ActiveTheme.Bg),
			LabelStyle: style.BgStyle().Foreground(style.ExtraDimColor).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Bottom),
		}
		imgViewer = boxLabelRenderer.Render(label, imgString, imgWidth)
	}

	displayHeight := 3
	renderVP := viewport.New(viewport.WithWidth(m.panelW-2), viewport.WithHeight(m.panelH-displayHeight-2))

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
