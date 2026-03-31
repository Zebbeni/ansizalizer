package modal

import (
	"path/filepath"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

// --- Export Current File ---

func (m Model) updateExportFile(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == FocusDestBrowser {
		var cmd tea.Cmd
		m.destBrowser, cmd = m.destBrowser.Update(msg)
		if m.destBrowser.ShouldClose {
			m.destBrowser.ShouldClose = false
			if m.destBrowser.SelectedDir != "" {
				m.focus = FocusConfirm
			}
		}
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		if m.focus == FocusConfirm {
			m.result = &Result{
				Kind:            ExportFileKind,
				SourcePath:      m.sourceFile,
				DestinationPath: m.destBrowser.SelectedDir,
				IsDir:           false,
			}
			m.Open = false
		} else if m.focus == FocusCancel {
			m.Open = false
		}
	case key.Matches(keyMsg, event.KeyMap.Left):
		if m.focus == FocusCancel {
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Right):
		if m.focus == FocusConfirm {
			m.focus = FocusCancel
		}
	case key.Matches(keyMsg, event.KeyMap.Up):
		if m.focus == FocusConfirm || m.focus == FocusCancel {
			m.focus = FocusDestBrowser
			m.destBrowser.SetActive(true)
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.Open = false
	}
	return m, nil
}

func (m Model) exportFileView() string {
	sourceLabel := style.DimmedTitle.Copy().Render("Source: " + filepath.Base(m.sourceFile))

	destLabel := style.DimmedTitle.Copy().Italic(true).
		Width(Width - 6).
		AlignHorizontal(lipgloss.Center).
		Render("Destination: " + filepath.Base(m.destBrowser.SelectedDir) + "/")

	if m.focus == FocusDestBrowser {
		return lipgloss.JoinVertical(lipgloss.Left, sourceLabel, "", destLabel, m.destBrowser.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, sourceLabel, "", destLabel)
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
			m.focus = FocusSourcePanel
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
			m.focus = FocusDestPanel
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
			m.focus = FocusSourceBrowser
			m.sourceBrowser.SetActive(true)
		case FocusDestPanel:
			m.focus = FocusDestBrowser
			m.destBrowser.SetActive(true)
		case FocusSubDirsYes:
			m.useSubDirs = true
		case FocusSubDirsNo:
			m.useSubDirs = false
		case FocusConfirm:
			m.result = &Result{
				Kind:            ExportBatchKind,
				SourcePath:      m.sourceBrowser.SelectedDir,
				DestinationPath: m.destBrowser.SelectedDir,
				IsDir:           true,
				UseSubDirs:      m.useSubDirs,
			}
			m.Open = false
		case FocusCancel:
			m.Open = false
		}

	case key.Matches(keyMsg, event.KeyMap.Right), key.Matches(keyMsg, event.KeyMap.Tab):
		switch m.focus {
		case FocusSourcePanel:
			m.focus = FocusDestPanel
		case FocusSubDirsYes:
			m.focus = FocusSubDirsNo
		case FocusConfirm:
			m.focus = FocusCancel
		}

	case key.Matches(keyMsg, event.KeyMap.Left):
		switch m.focus {
		case FocusDestPanel:
			m.focus = FocusSourcePanel
		case FocusSubDirsNo:
			m.focus = FocusSubDirsYes
		case FocusCancel:
			m.focus = FocusConfirm
		}

	case key.Matches(keyMsg, event.KeyMap.Down):
		switch m.focus {
		case FocusSourcePanel, FocusDestPanel:
			m.focus = FocusSubDirsYes
		case FocusSubDirsYes, FocusSubDirsNo:
			m.focus = FocusConfirm
		}

	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusSubDirsYes, FocusSubDirsNo:
			m.focus = FocusSourcePanel
		case FocusConfirm, FocusCancel:
			m.focus = FocusSubDirsYes
		}

	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.Open = false
	}
	return m, nil
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
	srcActive := m.focus == FocusSourceBrowser
	srcBox := m.renderPanel("Source", m.sourceBrowser.View(), srcFocused, srcActive, panelW, panelH)
	srcDir := style.DimmedTitle.Copy().Italic(true).Width(panelW + 2).AlignHorizontal(lipgloss.Center).
		Render(filepath.Base(m.sourceBrowser.SelectedDir) + "/")
	srcPanel := lipgloss.JoinVertical(lipgloss.Center, srcBox, srcDir)

	// Destination panel
	destFocused := m.focus == FocusDestPanel
	destActive := m.focus == FocusDestBrowser
	destBox := m.renderPanel("Destination", m.destBrowser.View(), destFocused, destActive, panelW, panelH)
	destDir := style.DimmedTitle.Copy().Italic(true).Width(panelW + 2).AlignHorizontal(lipgloss.Center).
		Render(filepath.Base(m.destBrowser.SelectedDir) + "/")
	destPanel := lipgloss.JoinVertical(lipgloss.Center, destBox, destDir)

	// Side by side panels with 1 char left margin
	panelRow := lipgloss.JoinHorizontal(lipgloss.Top, srcPanel, "  ", destPanel)
	panels := lipgloss.NewStyle().PaddingLeft(1).Render(panelRow)

	// Subdirectories toggle
	subTitle := style.DimmedTitle.Copy().Render("Include Subdirectories:")
	yesStyle := style.NormalButtonNode
	if m.focus == FocusSubDirsYes {
		yesStyle = style.FocusButtonNode
	} else if m.useSubDirs {
		yesStyle = style.ActiveButtonNode
	}
	noStyle := style.NormalButtonNode
	if m.focus == FocusSubDirsNo {
		noStyle = style.FocusButtonNode
	} else if !m.useSubDirs {
		noStyle = style.ActiveButtonNode
	}
	subRow := lipgloss.JoinHorizontal(lipgloss.Left,
		subTitle,
		style.BgStyle().PaddingLeft(1).Render(yesStyle.Render("Yes")),
		style.BgStyle().PaddingLeft(1).Render(noStyle.Render("No")),
	)

	// Description
	descStyle := style.DimmedTitle.Copy().Italic(true).
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
		subHelp := style.DimmedTitle.Copy().Italic(true).
			Render("Source directory structure will be mirrored in the destination.")
		parts = append(parts, subHelp)
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}
