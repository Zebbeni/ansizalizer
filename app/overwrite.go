package app

import (
	"path/filepath"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type OverwriteChoice int

const (
	OverwriteUndecided OverwriteChoice = iota
	OverwriteYes
	OverwriteSkip
	OverwriteCancel
)

type OverwriteFocus int

const (
	OverwriteFocusOverwrite OverwriteFocus = iota
	OverwriteFocusSkip
	OverwriteFocusCancel
	OverwriteFocusApplyAll
)

type OverwriteModal struct {
	filename   string
	focus      OverwriteFocus
	applyToAll bool
	choice     OverwriteChoice
}

func newOverwriteModal(destPath string) *OverwriteModal {
	return &OverwriteModal{
		filename: filepath.Base(destPath),
		focus:    OverwriteFocusOverwrite,
	}
}

func (m *OverwriteModal) Update(msg tea.KeyMsg) {
	switch {
	case key.Matches(msg, event.KeyMap.Left):
		switch m.focus {
		case OverwriteFocusSkip:
			m.focus = OverwriteFocusOverwrite
		case OverwriteFocusCancel:
			m.focus = OverwriteFocusSkip
		}
	case key.Matches(msg, event.KeyMap.Right), key.Matches(msg, event.KeyMap.Tab):
		switch m.focus {
		case OverwriteFocusOverwrite:
			m.focus = OverwriteFocusSkip
		case OverwriteFocusSkip:
			m.focus = OverwriteFocusCancel
		}
	case key.Matches(msg, event.KeyMap.Down):
		switch m.focus {
		case OverwriteFocusOverwrite, OverwriteFocusSkip, OverwriteFocusCancel:
			m.focus = OverwriteFocusApplyAll
		}
	case key.Matches(msg, event.KeyMap.Up):
		switch m.focus {
		case OverwriteFocusApplyAll:
			m.focus = OverwriteFocusOverwrite
		}
	case key.Matches(msg, event.KeyMap.Enter):
		switch m.focus {
		case OverwriteFocusOverwrite:
			m.choice = OverwriteYes
		case OverwriteFocusSkip:
			m.choice = OverwriteSkip
		case OverwriteFocusCancel:
			m.choice = OverwriteCancel
		case OverwriteFocusApplyAll:
			m.applyToAll = !m.applyToAll
		}
	case key.Matches(msg, event.KeyMap.Esc):
		m.choice = OverwriteCancel
	}
}

// handleOverwriteChoice processes the user's overwrite decision and continues the export.
func (m Model) handleOverwriteChoice() (Model, tea.Cmd) {
	choice := m.overwriteModal.choice
	applyAll := m.overwriteModal.applyToAll
	m.overwriteModal = nil

	if applyAll {
		m.overwriteAllChoice = choice
	}

	switch choice {
	case OverwriteSkip:
		return m.skipCurrentExport()
	case OverwriteCancel:
		m.waitingOnExport = false
		m.exportQueue = nil
		return m, event.BuildDisplayCmd("export cancelled")
	}

	// OverwriteYes — mark approved and re-trigger the render
	m.overwriteApproved = true
	return m, event.StartRenderToExportCmd
}

func (m *OverwriteModal) View() string {
	title := style.SelectedTitle.Copy().Render("File Exists")

	fileLabel := style.NormalTitle.Copy().
		Width(44).
		AlignHorizontal(lipgloss.Center).
		Render(m.filename)

	prompt := style.DimmedTitle.Copy().
		Width(44).
		AlignHorizontal(lipgloss.Center).
		Render("already exists. Overwrite?")

	// Buttons
	overwriteStyle := style.NormalButton
	if m.focus == OverwriteFocusOverwrite {
		overwriteStyle = style.FocusButton
	}
	skipStyle := style.NormalButton
	if m.focus == OverwriteFocusSkip {
		skipStyle = style.FocusButton
	}
	cancelStyle := style.NormalButton
	if m.focus == OverwriteFocusCancel {
		cancelStyle = style.FocusButton
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Center,
		overwriteStyle.Render(" Overwrite "),
		"  ",
		skipStyle.Render("  Skip  "),
		"  ",
		cancelStyle.Render(" Cancel "),
	)

	// Apply to all checkbox
	checkbox := style.RenderCheckbox("Apply to all", m.applyToAll, m.focus == OverwriteFocusApplyAll)

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		fileLabel,
		prompt,
		"",
		buttons,
		"",
		checkbox,
	)

	border := style.BgStyle().
		Border(style.HeavyBorder()).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg).
		Padding(1, 2).
		AlignHorizontal(lipgloss.Center)

	return style.ApplyBg(border.Render(content), 0)
}
