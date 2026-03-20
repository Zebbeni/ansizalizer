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

	displayMsg := fmt.Sprintf("viewing %s/%s with %s palette and %s", filepath.Base(filepath.Dir(msg.FilePath)), filepath.Base(msg.FilePath), msg.ColorsString, msg.AlphaString)
	return m, event.BuildDisplayCmd(displayMsg)
}
