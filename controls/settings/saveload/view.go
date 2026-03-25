package saveload

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	tabNames = map[Tab]string{
		LoadTab: "Import",
		SaveTab: "Export",
	}

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = style.BgStyle().Border(inactiveTabBorder, true)
	activeTabStyle    = style.BgStyle().Border(activeTabBorder, true)
	windowStyle       = style.BgStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0)
)

func (m Model) drawTabs() string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []Tab{LoadTab, SaveTab}
	tabFocus := map[Tab]State{SaveTab: SaveTabSelect, LoadTab: LoadTabSelect}

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	for i, t := range tabs {
		isFirst := i == 0
		isLast := i == len(tabs)-1
		isFocused := m.focus == tabFocus[t]
		isActiveTab := m.tab == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isFocused {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if isActiveTab {
				fgColor = style.NormalColor2
			}
		}

		var tabStyle lipgloss.Style
		if isActiveTab {
			tabStyle = style.ActiveTabStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if isFirst && isActiveTab {
			border.BottomLeft = "│"
		} else if isFirst && !isActiveTab {
			border.BottomLeft = "├"
		} else if isLast && isActiveTab {
			border.BottomRight = "└"
		} else if isLast && !isActiveTab {
			border.BottomRight = "┴"
		}

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Foreground(fgColor).Padding(0, 1)
		renderedTabs = append(renderedTabs, tabStyle.Render(tabNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW := max(m.width-lipgloss.Width(tabBlock)-2, 0)

	border := lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
	extendedStyle := style.TabWindowStyle.Copy().Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
	extended := extendedStyle.Copy().Width(extW).Height(1).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)

	return doc.String()
}

func (m Model) drawContent() string {
	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	tabRow := m.drawTabs()

	var content string
	switch m.tab {
	case SaveTab:
		content = m.drawSaveContent()
	case LoadTab:
		content = m.drawLoadContent()
	}

	body := style.TabWindowStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Width(lipgloss.Width(tabRow) - style.TabWindowStyle.GetHorizontalFrameSize()).Render(content)
	return tabRow + "\n" + body
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
		m.filenameInput.PromptStyle = nodeStyle
		m.filenameInput.TextStyle = textStyle
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
		parts = append(parts, style.BgStyle().Width(m.width-4).AlignHorizontal(lipgloss.Center).Render(browseText))
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
