package controls

import (
	"path/filepath"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawButtons() string {
	buttonWidth := m.width / numButtons
	buttons := make([]string, len(stateOrder))
	for i, state := range stateOrder {
		buttonStyle := style.NormalButton
		if m.active != Menu {
			// A tab is handling events
			if state == m.active {
				buttonStyle = style.ActiveButton
			} else if state == m.focus {
				buttonStyle = style.FocusButton
			}
		} else {
			// Tab bar has focus; show which tab's content is displayed
			if state == m.showing && m.showing != Menu {
				buttonStyle = style.ActiveButton
			}
			if state == m.focus {
				buttonStyle = style.FocusButton
			}
		}
		buttons[i] = buttonStyle.Copy().Width(buttonWidth).AlignHorizontal(lipgloss.Center).Render(stateNames[state])
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	return buttonRow
}

func (m Model) drawBrowserTitle() string {
	dir := filepath.Base(m.FileBrowser.SelectedDir)
	title := style.DimmedTitle.Copy().Italic(true).Render(dir + "/")
	return style.BgStyle().Width(m.width).Padding(1, 0).AlignHorizontal(lipgloss.Center).Render(title)
}
