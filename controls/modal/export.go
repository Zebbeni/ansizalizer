package modal

import (
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

// --- Export Current File ---

func (m Model) updateExportFile(msg tea.Msg) (Model, tea.Cmd) {
	// Browser active
	if m.focus == FocusDestBrowser {
		var cmd tea.Cmd
		m.destBrowser, cmd = m.destBrowser.Update(msg)
		if m.destBrowser.ShouldClose {
			m.destBrowser.ShouldClose = false
			m.destBrowser.SetActive(false)
			m.focus = FocusDestSelect
		}
		m.confirmedDest = m.destBrowser.SelectedDir
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		switch m.focus {
		case FocusDestPanel:
			m.focus = FocusDestSelect
		case FocusDestSelect:
			m.confirmedDest = m.destBrowser.SelectedDir
			m.focus = FocusConfirm
		case FocusConfirm:
			m.result = &Result{
				Kind:            ExportFileKind,
				SourcePath:      m.sourceFile,
				DestinationPath: m.confirmedDest,
				IsDir:           false,
			}
			m.Open = false
		case FocusCancel:
			m.Open = false
		}
	case key.Matches(keyMsg, event.KeyMap.Down), key.Matches(keyMsg, event.KeyMap.Tab):
		switch m.focus {
		case FocusDestPanel:
			m.focus = FocusConfirm
		case FocusDestSelect:
			m.focus = FocusDestBrowser
			m.destBrowser.SetActive(true)
		case FocusConfirm:
			m.focus = FocusCancel
		}
	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusDestSelect:
			m.focus = FocusDestPanel
		case FocusConfirm:
			m.focus = FocusDestPanel
		case FocusCancel:
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		switch m.focus {
		case FocusDestSelect:
			m.focus = FocusDestPanel
		default:
			m.Open = false
		}
	}
	return m, nil
}

func (m Model) exportFileView() string {
	contentArea := LoadWidth - 6
	panelContent := contentArea - 2

	sourceLabel := style.DimmedTitle.Italic(true).
		Width(contentArea).
		AlignHorizontal(lipgloss.Center).
		Render("Source: " + filepath.Base(m.sourceFile))

	// Destination panel — collapsed or expanded
	panelExpanded := m.focus == FocusDestSelect || m.focus == FocusDestBrowser
	panelFocused := m.focus == FocusDestPanel
	panelActive := m.focus == FocusDestBrowser || m.focus == FocusDestSelect

	renderColor := style.DimmedColor1
	if panelActive {
		renderColor = style.NormalColor1
	} else if panelFocused {
		renderColor = style.SelectedColor1
	}

	border := lipgloss.RoundedBorder()
	if panelFocused {
		border = style.HeavyBorder()
	}
	borderStyle := style.BgStyle().
		Border(border).
		BorderForeground(renderColor).
		BorderBackground(style.ActiveTheme.Bg)
	titleSt := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1).
		Foreground(renderColor)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: titleSt,
	}

	var panel string
	if panelExpanded {
		selectBtn := m.renderSelectButton(m.destBrowser.SelectedDir, m.focus == FocusDestSelect)
		browseContent := lipgloss.JoinVertical(lipgloss.Left, selectBtn, m.destBrowser.View())
		panel = renderer.Render("Destination", browseContent, panelContent)
	} else {
		summary := style.DimmedTitle.Italic(true).
			Width(panelContent).
			AlignHorizontal(lipgloss.Center).
			Padding(0, 1).
			Render(filepath.Base(m.confirmedDest) + "/")
		panel = renderer.Render("Destination", summary, panelContent)
	}

	// GIF explanation
	isGif := strings.HasSuffix(strings.ToLower(m.sourceFile), ".gif")
	if isGif {
		nameWithoutExt := strings.TrimSuffix(filepath.Base(m.sourceFile), filepath.Ext(m.sourceFile))
		gifNote := style.DimmedTitle.Italic(true).
			Width(contentArea).
			AlignHorizontal(lipgloss.Center).
			Render("Frames will be exported to numbered .ansi files in a new folder named " + nameWithoutExt + "/")
		return lipgloss.JoinVertical(lipgloss.Center, sourceLabel, "", panel, "", gifNote)
	}

	return lipgloss.JoinVertical(lipgloss.Center, sourceLabel, "", panel)
}

// exportFileName returns the output name shown on the Export button.
func (m Model) exportFileName() string {
	nameWithoutExt := strings.TrimSuffix(filepath.Base(m.sourceFile), filepath.Ext(m.sourceFile))
	isGif := strings.HasSuffix(strings.ToLower(m.sourceFile), ".gif")
	if isGif {
		return nameWithoutExt + "/"
	}
	return nameWithoutExt + ".ansi"
}

// --- Export Batch ---

func (m Model) updateExportBatch(msg tea.Msg) (Model, tea.Cmd) {
	// Source browser active — route all input to it
	if m.focus == FocusSourceBrowser {
		var cmd tea.Cmd
		m.sourceBrowser, cmd = m.sourceBrowser.Update(msg)
		if m.sourceBrowser.ShouldClose {
			m.sourceBrowser.ShouldClose = false
			m.sourceBrowser.SetActive(false)
			m.focus = FocusSourceSelect
		}
		return m, cmd
	}

	// Destination browser active
	if m.focus == FocusDestBrowser {
		var cmd tea.Cmd
		m.destBrowser, cmd = m.destBrowser.Update(msg)
		if m.destBrowser.ShouldClose {
			m.destBrowser.ShouldClose = false
			m.destBrowser.SetActive(false)
			m.focus = FocusDestSelect
		}
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		switch m.focus {
		case FocusSourcePanel:
			m.focus = FocusSourceSelect
		case FocusSourceSelect:
			// Confirm source directory and tab to next item
			m.confirmedSrc = m.sourceBrowser.SelectedDir
			m.imageCount = countImages(m.confirmedSrc, m.useSubDirs)
			m.focus = FocusDestPanel
		case FocusDestPanel:
			m.focus = FocusDestSelect
		case FocusDestSelect:
			// Confirm dest directory and tab to next item
			m.confirmedDest = m.destBrowser.SelectedDir
			m.focus = FocusSubDirsYes
		case FocusSubDirsYes, FocusSubDirsNo:
			m.useSubDirs = !m.useSubDirs
			if m.confirmedSrc != "" {
				m.imageCount = countImages(m.confirmedSrc, m.useSubDirs)
			}
		case FocusConfirm:
			if m.confirmedSrc != "" && m.confirmedDest != "" {
				m.result = &Result{
					Kind:            ExportBatchKind,
					SourcePath:      m.confirmedSrc,
					DestinationPath: m.confirmedDest,
					IsDir:           true,
					UseSubDirs:      m.useSubDirs,
				}
				m.Open = false
			}
		case FocusCancel:
			m.Open = false
		}

	case key.Matches(keyMsg, event.KeyMap.Right), key.Matches(keyMsg, event.KeyMap.Tab):
		switch m.focus {
		case FocusSourcePanel:
			m.focus = FocusDestPanel
		case FocusSourceSelect:
			m.focus = FocusDestPanel
		case FocusDestSelect:
			m.focus = FocusSubDirsYes
		case FocusSubDirsYes, FocusSubDirsNo:
			m.focus = FocusConfirm
		case FocusConfirm:
			m.focus = FocusCancel
		}

	case key.Matches(keyMsg, event.KeyMap.Left):
		switch m.focus {
		case FocusDestPanel:
			m.focus = FocusSourcePanel
		case FocusDestSelect:
			m.focus = FocusSourcePanel
		case FocusConfirm:
			m.focus = FocusSubDirsYes
		case FocusCancel:
			m.focus = FocusConfirm
		}

	case key.Matches(keyMsg, event.KeyMap.Down):
		switch m.focus {
		case FocusSourcePanel, FocusDestPanel:
			m.focus = FocusSubDirsYes
		case FocusSourceSelect:
			m.focus = FocusSourceBrowser
			m.sourceBrowser.SetActive(true)
		case FocusDestSelect:
			m.focus = FocusDestBrowser
			m.destBrowser.SetActive(true)
		case FocusSubDirsYes, FocusSubDirsNo:
			m.focus = FocusConfirm
		}

	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusSourceSelect:
			m.focus = FocusSourcePanel
		case FocusDestSelect:
			m.focus = FocusDestPanel
		case FocusSubDirsYes, FocusSubDirsNo:
			m.focus = FocusSourcePanel
		case FocusConfirm, FocusCancel:
			m.focus = FocusSubDirsYes
		}

	case key.Matches(keyMsg, event.KeyMap.Esc):
		switch m.focus {
		case FocusSourceSelect:
			m.focus = FocusSourcePanel
		case FocusDestSelect:
			m.focus = FocusDestPanel
		default:
			m.Open = false
		}
	}
	return m, nil
}

func (m Model) renderSelectButton(dir string, focused bool) string {
	label := "Select " + filepath.Base(dir) + " \U0001F5C0"
	btnStyle := style.DimmedTitle
	if focused {
		btnStyle = style.SelectedTitle
	}
	return btnStyle.Padding(1, 0, 1, 2).Render(label)
}

func (m Model) renderPanel(title, content string, focused, active bool, w, h int) string {
	renderColor := style.DimmedColor1
	if active {
		renderColor = style.NormalColor1
	} else if focused {
		renderColor = style.SelectedColor1
	}

	textStyle := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1).
		Foreground(renderColor)
	border := lipgloss.RoundedBorder()
	if focused && !active {
		border = style.HeavyBorder()
	}
	borderStyle := style.BgStyle().
		Border(border).
		BorderForeground(renderColor).
		BorderBackground(style.ActiveTheme.Bg).
		Height(h)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: textStyle,
	}

	return renderer.Render(title, content, w)
}

func (m Model) exportBatchView() string {
	// BoxWithLabel.Render(w-2) produces total width w.
	// Inner content = w - 2 (border). Two panels + 2 gap, each panel adds 2 border.
	// (panelW + 2) + 2 + (panelW + 2) = w - 2 → panelW = (w - 8) / 2
	panelW := (m.width - 9) / 2

	panelH := BatchHeight - 14

	// Source panel
	srcFocused := m.focus == FocusSourcePanel
	srcActive := m.focus == FocusSourceBrowser || m.focus == FocusSourceSelect
	srcSelectBtn := m.renderSelectButton(m.sourceBrowser.SelectedDir, m.focus == FocusSourceSelect)
	srcContent := lipgloss.JoinVertical(lipgloss.Left, srcSelectBtn, m.sourceBrowser.View())
	srcBox := m.renderPanel("Source", srcContent, srcFocused, srcActive, panelW, panelH)
	srcConfirmed := m.confirmedSrc
	if srcConfirmed == "" {
		srcConfirmed = "Source: (not selected)"
	} else {
		srcConfirmed = "Source: " + filepath.Base(srcConfirmed) + "/"
	}
	srcDir := style.DimmedTitle.Italic(true).Width(panelW + 2).AlignHorizontal(lipgloss.Center).
		Render(srcConfirmed)
	srcPanel := lipgloss.JoinVertical(lipgloss.Center, srcBox, srcDir)

	// Destination panel
	destFocused := m.focus == FocusDestPanel
	destActive := m.focus == FocusDestBrowser || m.focus == FocusDestSelect
	destSelectBtn := m.renderSelectButton(m.destBrowser.SelectedDir, m.focus == FocusDestSelect)
	destContent := lipgloss.JoinVertical(lipgloss.Left, destSelectBtn, m.destBrowser.View())
	destBox := m.renderPanel("Destination", destContent, destFocused, destActive, panelW, panelH)
	destConfirmed := m.confirmedDest
	if destConfirmed == "" {
		destConfirmed = "Destination: (not selected)"
	} else {
		destConfirmed = "Destination: " + filepath.Base(destConfirmed) + "/"
	}
	destDir := style.DimmedTitle.Italic(true).Width(panelW + 2).AlignHorizontal(lipgloss.Center).
		Render(destConfirmed)
	destPanel := lipgloss.JoinVertical(lipgloss.Center, destBox, destDir)

	// Side by side panels with 1 char left margin
	panelRow := lipgloss.JoinHorizontal(lipgloss.Top, srcPanel, "  ", destPanel)
	panels := lipgloss.NewStyle().PaddingLeft(1).Render(panelRow)

	// Subdirectories toggle
	subFocused := m.focus == FocusSubDirsYes || m.focus == FocusSubDirsNo
	subRow := style.RenderCheckbox("Include Subdirectories", m.useSubDirs, subFocused)

	// Description
	descStyle := style.DimmedTitle.Italic(true).
		Width(m.width - 10).
		AlignHorizontal(lipgloss.Center)
	desc1 := descStyle.Render("Render all images in the source directory as .ansi files using the current settings.")
	desc2 := descStyle.Render("Animated GIFs produce folders with numbered frame files.")
	descRendered := lipgloss.NewStyle().Padding(1, 0).Render(
		lipgloss.JoinVertical(lipgloss.Center, desc1, desc2))

	subRow = lipgloss.NewStyle().PaddingTop(1).Render(subRow)

	parts := []string{descRendered, panels, subRow}

	// Show subdirectory help near the toggle
	if m.useSubDirs {
		subHelp := style.DimmedTitle.Italic(true).
			Render("Source directory structure will be mirrored in the destination.")
		parts = append(parts, subHelp)
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}
