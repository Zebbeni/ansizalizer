package saveload

import (
	"fmt"
	"path/filepath"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var tabNames = map[Tab]string{
	LoadTab: "Import",
	SaveTab: "Export",
}

func (m Model) drawContent() string {
	tabs := []style.Tab{
		{Label: tabNames[LoadTab], Focused: m.focus == LoadTabSelect, Active: m.tab == LoadTab},
		{Label: tabNames[SaveTab], Focused: m.focus == SaveTabSelect, Active: m.tab == SaveTab},
	}
	var content string
	switch m.tab {
	case SaveTab:
		content = m.drawSaveContent()
	case LoadTab:
		content = m.drawLoadContent()
	}
	return style.RenderTabs(tabs, m.IsActive, content, m.width)
}

func (m Model) drawSaveContent() string {
	centeredStyle := style.BgStyle().Width(m.width - 4).AlignHorizontal(lipgloss.Center)
	parts := make([]string, 0, 5)

	// Save result message (shown after save attempt)
	if m.saveResultText != "" {
		result := centeredStyle.Render(style.DimmedTitle.Copy().Render(m.saveResultText))
		parts = append(parts, result)
	}

	// Directory browse prompt or selected directory
	inputValueStyle := style.DimmedTitle.Copy()
	if m.focus == DirInput && m.IsActive {
		inputValueStyle = style.SelectedTitle.Copy()
	} else if m.focus == DirInput {
		inputValueStyle = style.NormalTitle.Copy()
	}

	if m.dirChosen {
		dirTitle := centeredStyle.Render(style.DimmedTitle.Copy().Render("Directory"))
		parent := filepath.Base(filepath.Dir(m.selectedDir))
		selected := filepath.Base(m.selectedDir)
		dirValue := fmt.Sprintf("%s/%s", parent, selected)
		dirRunes := []rune(dirValue)
		if len(dirRunes) > m.width-4 {
			dirValue = ".." + string(dirRunes[len(dirRunes)-(m.width-6):])
		}
		dirContent := lipgloss.JoinVertical(lipgloss.Left, dirTitle, centeredStyle.Render(inputValueStyle.Render(dirValue)))
		parts = append(parts, dirContent)
	} else {
		browseText := inputValueStyle.Render("Choose Directory ↴")
		parts = append(parts, centeredStyle.Render(browseText))
	}

	// Directory browser (only shown when active)
	if m.focus == DirBrowser {
		dir := filepath.Base(m.dirBrowser.SelectedDir)
		title := style.DimmedTitle.Copy().Italic(true).Padding(1, 0).Render("Browsing " + dir + "/")
		parts = append(parts, centeredStyle.Render(title))
		parts = append(parts, m.dirBrowser.View())
	}

	// Filename input (only shown after directory chosen)
	if m.dirChosen {
		nodeStyle := style.NormalButtonNode.Copy().PaddingRight(1)
		textStyle := style.BgStyle().Foreground(style.DimmedColor1)
		if m.filenameInput.Focused() {
			nodeStyle = style.ActiveButtonNode.Copy().PaddingRight(1)
			textStyle = style.BgStyle().Foreground(style.SelectedColor1)
		} else if m.focus == FilenameForm && m.IsActive {
			nodeStyle = style.FocusButtonNode.Copy().PaddingRight(1)
			textStyle = style.BgStyle().Foreground(style.NormalColor1)
		}
		fnStyles := m.filenameInput.Styles()
		fnStyles.Focused.Prompt = nodeStyle
		fnStyles.Focused.Text = textStyle
		fnStyles.Blurred.Prompt = nodeStyle
		fnStyles.Blurred.Text = textStyle
		fnStyles.Cursor.Blink = m.filenameInput.Focused()
		m.filenameInput.SetStyles(fnStyles)
		m.filenameInput.SetVirtualCursor(m.filenameInput.Focused())
		parts = append(parts, m.filenameInput.View())
	}

	// Save/Cancel buttons (only shown after filename entered)
	if m.dirChosen && m.filenameInput.Value() != "" && !m.filenameInput.Focused() && m.focus != DirBrowser {
		filename := m.filenameInput.Value() + ".json"
		prompt := centeredStyle.Copy().PaddingTop(1).Render(style.DimmedTitle.Copy().Render(fmt.Sprintf("Save %s?", filename)))
		parts = append(parts, prompt)

		saveStyle := style.NormalButton
		if m.focus == SaveButton && m.IsActive {
			saveStyle = style.FocusButton
		}
		cancelStyle := style.NormalButton
		if m.focus == CancelButton && m.IsActive {
			cancelStyle = style.FocusButton
		}
		buttonRow := lipgloss.JoinHorizontal(lipgloss.Center, saveStyle.Render("  Save  "), "  ", cancelStyle.Render(" Cancel "))
		buttons := centeredStyle.Render(buttonRow)
		parts = append(parts, buttons)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) drawLoadContent() string {
	parts := make([]string, 0, 4)

	// Load result message (shown after load attempt)
	if m.loadResultText != "" {
		resultStyle := style.DimmedTitle.Copy()
		result := style.BgStyle().Width(m.width - 4).AlignHorizontal(lipgloss.Center).Render(resultStyle.Render(m.loadResultText))
		parts = append(parts, result)
	}

	// Browse input / confirmation
	inputValueStyle := style.DimmedTitle.Copy()
	if m.focus == LoadInput && m.IsActive {
		inputValueStyle = style.SelectedTitle.Copy()
	} else if m.focus == LoadInput {
		inputValueStyle = style.NormalTitle.Copy()
	}

	if m.selectedLoadFile != "" {
		// Show the selected path with confirmation
		parent := filepath.Base(filepath.Dir(m.selectedLoadFile))
		selected := filepath.Base(m.selectedLoadFile)
		fileValue := fmt.Sprintf("%s/%s", parent, selected)
		fileRunes := []rune(fileValue)
		if len(fileRunes) > m.width-4 {
			fileValue = ".." + string(fileRunes[len(fileRunes)-(m.width-6):])
		}

		filename := filepath.Base(m.selectedLoadFile)
		prompt := style.DimmedTitle.Copy().Render(fmt.Sprintf("Load %s?", filename))
		parts = append(parts, style.BgStyle().Width(m.width-4).AlignHorizontal(lipgloss.Center).PaddingTop(1).Render(prompt))

		loadStyle := style.NormalButton
		if m.focus == LoadConfirmButton && m.IsActive {
			loadStyle = style.FocusButton
		}
		cancelStyle := style.NormalButton
		if m.focus == CancelButton && m.IsActive {
			cancelStyle = style.FocusButton
		}
		buttonRow := lipgloss.JoinHorizontal(lipgloss.Center, loadStyle.Render("  Load  "), "  ", cancelStyle.Render(" Cancel "))
		buttons := style.BgStyle().Width(m.width - 4).AlignHorizontal(lipgloss.Center).Render(buttonRow)
		parts = append(parts, buttons)
	} else {
		// Show browse prompt
		browseText := inputValueStyle.Render("Browse Settings Files ↴")
		parts = append(parts, style.BgStyle().Width(m.width-2).AlignHorizontal(lipgloss.Center).Render(browseText))
	}

	// File browser (only shown when active)
	if m.focus == LoadBrowser {
		dir := filepath.Base(m.loadBrowser.SelectedDir)
		centeredStyle := style.BgStyle().Width(m.width - 4).AlignHorizontal(lipgloss.Center)
		title := style.DimmedTitle.Copy().Italic(true).Padding(1, 0).Render("Browsing " + dir + "/")
		parts = append(parts, centeredStyle.Render(title))
		parts = append(parts, m.loadBrowser.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
