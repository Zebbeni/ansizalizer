package source

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
	valueStyle := style.DimmedTitle
	if SrcInput == m.focus {
		valueStyle = style.SelectedTitle
	}

	value := m.SourceBrowser.ActiveFile
	if m.doExportDirectory {
		value = m.SourceBrowser.ActiveDir
	}
	value = filepath.Base(value)

	if value == "." || len(value) == 0 {
		value = "(None)"
	}

	valueContent := valueStyle.Render(value)

	valueWidth := m.width
	widthStyle := lipgloss.NewStyle().Width(valueWidth).AlignHorizontal(lipgloss.Center)
	valueContent = widthStyle.Render(valueContent)

	if m.focus != SrcBrowser {
		return valueContent
	}

	browserContent := m.SourceBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Left, valueContent, browserContent)
}
