package viewer

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleFinishRenderMsg(msg event.FinishRenderToViewMsg) (Model, tea.Cmd) {
	m.WaitingOnRender = false
	m.imgString = msg.ImgString

	// Clear animation state
	m.frames = nil
	m.currentFrame = 0
	m.isAnimating = false

	displayMsg := fmt.Sprintf("viewing %s/%s with %s palette and %s", filepath.Base(filepath.Dir(msg.FilePath)), filepath.Base(msg.FilePath), msg.ColorsString, msg.AlphaString)
	return m, event.BuildDisplayCmd(displayMsg)
}

func (m Model) handleFinishRenderGIFMsg(msg event.FinishRenderGIFToViewMsg) (Model, tea.Cmd) {
	m.WaitingOnRender = false
	m.frames = msg.Frames
	m.delay = msg.Delay
	m.currentFrame = 0
	m.isAnimating = true
	m.imgString = ""

	displayMsg := fmt.Sprintf("viewing %s/%s (%d frames, %dms) with %s palette and %s",
		filepath.Base(filepath.Dir(msg.FilePath)),
		filepath.Base(msg.FilePath),
		len(msg.Frames),
		msg.Delay.Milliseconds(),
		msg.ColorsString,
		msg.AlphaString,
	)

	return m, tea.Batch(
		event.BuildDisplayCmd(displayMsg),
		event.BuildAnimationTickCmd(m.delay),
	)
}

func (m Model) handleAnimationTick() (Model, tea.Cmd) {
	if !m.isAnimating || len(m.frames) == 0 || m.WaitingOnRender {
		m.isAnimating = false
		return m, nil
	}

	m.currentFrame = (m.currentFrame + 1) % len(m.frames)
	return m, event.BuildAnimationTickCmd(m.delay)
}
