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
	msg := style.DimmedTitle.Copy().Italic(true).
		AlignHorizontal(lipgloss.Center).
		Render("Are you sure you want to close the app?")

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

	content := lipgloss.JoinVertical(lipgloss.Center, msg, "", buttons)

	textStyle := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1).
		Foreground(style.SelectedColor1)
	borderStyle := style.BgStyle().
		Border(style.HeavyBorder()).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg).
		Padding(1, 2)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: textStyle,
	}

	return style.ApplyBg(renderer.Render("Quit", content, lipgloss.Width(content)), 0)
}
