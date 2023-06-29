package export

import (
	"path/filepath"

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
	widthStyle := lipgloss.NewStyle().Width(m.width / 2)
	optionStyle := style.NormalButtonNode
	if ExpFile == m.focus {
		optionStyle = style.FocusButtonNode
	} else if m.doExportDirectory == false {
		optionStyle = style.ActiveButtonNode
	}
	singleFile := optionStyle.Copy().Render(stateNames[ExpFile])
	singleFile = widthStyle.Render(singleFile)

	optionStyle = style.NormalButtonNode
	if ExpDirectory == m.focus {
		optionStyle = style.FocusButtonNode
	} else if m.doExportDirectory {
		optionStyle = style.ActiveButtonNode
	}
	directory := optionStyle.Copy().Render(stateNames[ExpDirectory])
	directory = widthStyle.Render(directory)

	return lipgloss.JoinHorizontal(lipgloss.Center, singleFile, " ", directory)
}

func (m Model) drawSource() string {
	promptStyle := style.NormalButtonNode
	valueStyle := style.DimmedTitle
	if SrcInput == m.focus {
		promptStyle = style.FocusButtonNode
		valueStyle = style.SelectedTitle
	}
	promptString := promptStyle.Render("Source: ")

	value := m.SourceBrowser.ActiveFile
	if m.doExportDirectory {
		value = m.SourceBrowser.ActiveDir
	}
	value = filepath.Base(value)

	valueString := valueStyle.Render(value)

	valueWidth := m.width - lipgloss.Width(promptString)
	widthStyle := lipgloss.NewStyle().Width(valueWidth)
	valueString = widthStyle.Render(valueString)

	source := lipgloss.JoinHorizontal(lipgloss.Center, promptString, valueString)

	if m.focus != SrcBrowser {
		return source
	}

	browser := m.SourceBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Left, source, browser)
}
