package viewer

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleFinishRenderMsg(msg event.FinishRenderToViewMsg) (Model, tea.Cmd) {
	m.WaitingOnRender = false

	// Clear animation state
	m.frames = nil
	m.currentFrame = 0
	m.isAnimating = false

	m.imgString = msg.ImgString
	m.wrappedView = m.buildWrappedView()

	displayMsg := fmt.Sprintf("viewing %s/%s with %s", filepath.Base(filepath.Dir(msg.FilePath)), filepath.Base(msg.FilePath), msg.ColorsString)
	if msg.AlphaString != "" {
		displayMsg += ", " + msg.AlphaString
	}
	if msg.ExtraInfo != "" {
		displayMsg += ", " + msg.ExtraInfo
	}
	return m, event.BuildDisplayCmd(displayMsg)
}

func (m Model) handleFinishRenderGIFMsg(msg event.FinishRenderGIFToViewMsg) (Model, tea.Cmd) {
	m.WaitingOnRender = false
	m.frames = msg.Frames
	m.delay = msg.Delay
	m.currentFrame = 0
	m.isAnimating = true
	m.imgString = msg.Frames[0]
	m.generation++
	m.wrappedView = m.buildWrappedView()

	displayMsg := fmt.Sprintf("viewing %s/%s (%d frames, %dms) with %s",
		filepath.Base(filepath.Dir(msg.FilePath)),
		filepath.Base(msg.FilePath),
		len(msg.Frames),
		msg.Delay.Milliseconds(),
		msg.ColorsString,
	)
	if msg.AlphaString != "" {
		displayMsg += ", " + msg.AlphaString
	}
	if msg.ExtraInfo != "" {
		displayMsg += ", " + msg.ExtraInfo
	}

	return m, tea.Batch(
		event.BuildDisplayCmd(displayMsg),
		event.BuildAnimationTickCmd(m.delay, m.generation),
	)
}

func (m Model) handleAnimationTick(msg event.AnimationTickMsg) (Model, tea.Cmd) {
	if !m.isAnimating || len(m.frames) == 0 {
		return m, nil
	}

	// Ignore ticks from a previous generation
	if msg.Generation != m.generation {
		return m, nil
	}

	if m.WaitingOnRender {
		// Keep showing current frame while re-render is in progress
		return m, event.BuildAnimationTickCmd(m.delay, m.generation)
	}

	m.currentFrame = (m.currentFrame + 1) % len(m.frames)
	m.imgString = m.frames[m.currentFrame]
	m.wrappedView = m.buildWrappedView()

	return m, event.BuildAnimationTickCmd(m.delay, m.generation)
}
