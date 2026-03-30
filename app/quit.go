package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) handleQuitModalKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Enter):
		if m.quitFocusYes {
			return m, tea.Quit
		}
		m.quitConfirm = false
	case key.Matches(msg, event.KeyMap.Left):
		m.quitFocusYes = true
	case key.Matches(msg, event.KeyMap.Right):
		m.quitFocusYes = false
	case key.Matches(msg, event.KeyMap.Esc):
		m.quitConfirm = false
	case key.Matches(msg, event.KeyMap.Quit):
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) renderQuitModal() string {
	title := style.SelectedTitle.Copy().
		AlignHorizontal(lipgloss.Center).
		Width(20).
		Render("Quit Program?")

	yesStyle := style.NormalButton
	noStyle := style.NormalButton
	if m.quitFocusYes {
		yesStyle = style.FocusButton
	} else {
		noStyle = style.FocusButton
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center,
		yesStyle.Render("  Yes  "),
		"  ",
		noStyle.Render("  No   "),
	)

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", buttons)

	border := style.BgStyle().
		Padding(1, 2).
		Border(style.HeavyBorder()).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg).
		AlignHorizontal(lipgloss.Center)

	return style.ApplyBg(border.Render(content), 0)
}
