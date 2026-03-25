package viewer

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/event"
)

type Model struct {
	imgString string
	settings  settings.Model

	// Animation state
	frames       []string
	delay        time.Duration
	currentFrame int
	isAnimating  bool
	generation   int

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

func (m Model) IsAnimating() bool {
	return m.isAnimating && len(m.frames) > 1
}

func (m Model) Frames() []string {
	return m.frames
}

func (m Model) View() string {
	if m.isAnimating && len(m.frames) > 0 {
		return m.frames[m.currentFrame]
	}
	return m.imgString
}
