package source

import (
	"fmt"
	"path/filepath"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateNames = map[State]string{
		ExpFile:      "Single File",
		ExpDirectory: "Directory",
	}
)

func (m Model) drawExportTypeOptions() string {
	buttonWidth := (m.width / 2) - 2

	optionStyle := style.NormalButton
	if ExpFile == m.focus && m.IsActive {
		optionStyle = style.FocusButton
	} else if m.doExportDirectory == false {
		optionStyle = style.ActiveButton
	}
	singleFileButton := optionStyle.Width(buttonWidth).AlignHorizontal(lipgloss.Center).Render(stateNames[ExpFile])

	optionStyle = style.NormalButton
	if ExpDirectory == m.focus && m.IsActive {
		optionStyle = style.FocusButton
	} else if m.doExportDirectory {
		optionStyle = style.ActiveButton
	}
	directoryButton := optionStyle.Width(buttonWidth).AlignHorizontal(lipgloss.Center).Render(stateNames[ExpDirectory])

	return lipgloss.JoinHorizontal(lipgloss.Center, singleFileButton, directoryButton)
}

func (m Model) drawSubDirOptions() string {
	title := style.DimmedTitle.Render("Include Subdirectories")

	nodeWidthStyle := style.BgStyle().Width(m.width / 2).AlignHorizontal(lipgloss.Center)

	yesStyle := style.NormalButtonNode
	if m.includeSubdirectories {
		yesStyle = style.ActiveButtonNode
	}
	if m.focus == SubDirsYes && m.IsActive {
		yesStyle = style.FocusButtonNode
	}
	yesNode := nodeWidthStyle.Render(yesStyle.Render("Yes"))

	noStyle := style.NormalButtonNode
	if !m.includeSubdirectories {
		noStyle = style.ActiveButtonNode
	}
	if m.focus == SubDirsNo && m.IsActive {
		noStyle = style.FocusButtonNode
	}

	noStyle.Padding(0)
	noNode := nodeWidthStyle.Render(noStyle.Render("No"))

	options := lipgloss.JoinHorizontal(lipgloss.Center, yesNode, noNode)

	widthStyle := style.BgStyle().Width(m.width).AlignHorizontal(lipgloss.Left).PaddingBottom(1)
	content := lipgloss.JoinVertical(lipgloss.Center, title, options)

	return widthStyle.Render(content)
}

func (m Model) drawPrompt() string {
	return style.DimmedTitle.AlignHorizontal(lipgloss.Center).Padding(0).Render("Select")
}

func (m Model) drawSelected() string {
	title := style.DimmedTitle.Render("Selected")

	valueStyle := style.DimmedTitle
	if Input == m.focus {
		if m.IsActive {
			valueStyle = style.SelectedTitle
		} else {
			valueStyle = style.NormalTitle
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

	widthStyle := style.BgStyle().Width(m.width).AlignHorizontal(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, title, valueContent)

	return widthStyle.Render(content)
}

func (m Model) drawBrowserTitle() string {
	dir := filepath.Base(m.Browser.SelectedDir)
	if m.doExportDirectory {
		return style.DimmedTitle.Italic(true).Padding(0, 2, 1, 2).Render("Browsing " + dir + "/")
	}
	return style.DimmedTitle.Italic(true).Padding(0, 2, 1, 2).Render("Browsing " + dir + "/")
}
