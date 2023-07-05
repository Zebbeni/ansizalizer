package source

import (
	"fmt"
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
	widthStyle := lipgloss.NewStyle().Width((m.width / 2) - 2).AlignHorizontal(lipgloss.Center)
	optionStyle := style.NormalButton
	if ExpFile == m.focus && m.IsActive {
		optionStyle = style.FocusButton
	} else if m.doExportDirectory == false {
		optionStyle = style.ActiveButton
	}
	singleFileButtonText := widthStyle.Render(stateNames[ExpFile])
	singleFileButton := optionStyle.Render(singleFileButtonText)

	optionStyle = style.NormalButton
	if ExpDirectory == m.focus && m.IsActive {
		optionStyle = style.FocusButton
	} else if m.doExportDirectory {
		optionStyle = style.ActiveButton
	}
	directoryButtonText := widthStyle.Render(stateNames[ExpDirectory])
	directoryButton := optionStyle.Render(directoryButtonText)

	return lipgloss.JoinHorizontal(lipgloss.Center, singleFileButton, directoryButton)
}

func (m Model) drawSubDirOptions() string {
	title := style.DimmedTitle.Copy().Render("Include Subdirectories")

	nodeWidthStyle := lipgloss.NewStyle().Width(m.width / 2).AlignHorizontal(lipgloss.Center)

	yesStyle := style.NormalButtonNode.Copy()
	if m.includeSubdirectories {
		yesStyle = style.ActiveButtonNode.Copy()
	}
	if m.focus == SubDirsYes {
		yesStyle = style.FocusButtonNode.Copy()
	}
	yesNode := nodeWidthStyle.Render(yesStyle.Render("Yes"))

	noStyle := style.NormalButtonNode.Copy()
	if !m.includeSubdirectories {
		noStyle = style.ActiveButtonNode.Copy()
	}
	if m.focus == SubDirsNo {
		noStyle = style.FocusButtonNode.Copy()
	}

	noStyle.Padding(0)
	noNode := nodeWidthStyle.Render(noStyle.Render("No"))

	options := lipgloss.JoinHorizontal(lipgloss.Center, yesNode, noNode)

	widthStyle := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Left).PaddingBottom(1)
	content := lipgloss.JoinVertical(lipgloss.Center, title, options)

	return widthStyle.Render(content)
}

func (m Model) drawPrompt() string {
	return style.DimmedTitle.Copy().AlignHorizontal(lipgloss.Center).Padding(0).Render("Select")
}

func (m Model) drawSelected() string {
	title := style.DimmedTitle.Copy().Render("Selected")

	valueStyle := style.DimmedTitle.Copy()
	if Input == m.focus {
		if m.IsActive {
			valueStyle = style.SelectedTitle.Copy()
		} else {
			valueStyle = style.NormalTitle.Copy()
		}
	}
	valueStyle.Padding(0, 0, 1, 0)

	path := m.Browser.SelectedFile
	if m.doExportDirectory {
		path = m.Browser.SelectedDir
	}

	parent := filepath.Base(filepath.Dir(path))
	selected := filepath.Base(path)
	value := fmt.Sprintf("%s/%s", parent, selected)

	valueRunes := []rune(value)
	if len(valueRunes) > m.width {
		value = string(valueRunes[len(valueRunes)-m.width:])
	}

	valueContent := valueStyle.Render(value)

	widthStyle := lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, title, valueContent)

	return widthStyle.Render(content)
}
