package export

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateNames = map[State]string{
		ExpFile:      "Single File",
		ExpDirectory: "Directory",
	}
)

func (m Model) drawExportTypeOptions() string {
	optionStyle := style.NormalButtonNode
	if ExpFile == m.focus {
		optionStyle = style.FocusButtonNode
	} else if m.doExportDirectory == false {
		optionStyle = style.ActiveButtonNode
	}
	w := m.width / 2
	singleFile := optionStyle.Copy().Width(w).Render(stateNames[ExpFile])

	optionStyle = style.NormalButtonNode
	if ExpDirectory == m.focus {
		optionStyle = style.FocusButtonNode
	} else if m.doExportDirectory {
		optionStyle = style.ActiveButtonNode
	}
	directory := optionStyle.Copy().Width(w).Render(stateNames[ExpDirectory])

	return lipgloss.JoinHorizontal(lipgloss.Left, singleFile, directory)
}
