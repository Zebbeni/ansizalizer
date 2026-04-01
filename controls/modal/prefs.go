package modal

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) updatePrefs(msg tea.Msg) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		switch m.focus {
		case FocusPrefsShowHelp:
			m.prefsShowHelp = !m.prefsShowHelp
		case FocusPrefsRestoreTheme:
			m.prefsRestoreTheme = !m.prefsRestoreTheme
		case FocusPrefsRestoreSettings:
			m.prefsRestoreSettings = !m.prefsRestoreSettings
		case FocusPrefsClose:
			m.result = &Result{
				Kind:            PrefsKind,
				ShowHelp:        m.prefsShowHelp,
				RestoreTheme:    m.prefsRestoreTheme,
				RestoreSettings: m.prefsRestoreSettings,
			}
			m.Open = false
		}
	case key.Matches(keyMsg, event.KeyMap.Down), key.Matches(keyMsg, event.KeyMap.Tab):
		switch m.focus {
		case FocusPrefsShowHelp:
			m.focus = FocusPrefsRestoreTheme
		case FocusPrefsRestoreTheme:
			m.focus = FocusPrefsRestoreSettings
		case FocusPrefsRestoreSettings:
			m.focus = FocusPrefsClose
		}
	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusPrefsRestoreTheme:
			m.focus = FocusPrefsShowHelp
		case FocusPrefsRestoreSettings:
			m.focus = FocusPrefsRestoreTheme
		case FocusPrefsClose:
			m.focus = FocusPrefsRestoreSettings
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.result = &Result{
			Kind:            PrefsKind,
			ShowHelp:        m.prefsShowHelp,
			RestoreTheme:    m.prefsRestoreTheme,
			RestoreSettings: m.prefsRestoreSettings,
		}
		m.Open = false
	}
	return m, nil
}

func (m Model) prefsView() string {
	w := Width - 8

	desc := style.DimmedTitle.Italic(true).
		Width(w).
		AlignHorizontal(lipgloss.Center).
		Render("Settings applied on each app start.")

	labelWidth := 27 // wide enough for "Restore Previous Settings"
	showHelp := style.RenderCheckboxFixedWidth("Show Help",
		m.prefsShowHelp, m.focus == FocusPrefsShowHelp, labelWidth)

	restoreTheme := style.RenderCheckboxFixedWidth("Restore Previous Theme",
		m.prefsRestoreTheme, m.focus == FocusPrefsRestoreTheme, labelWidth)

	restoreSettings := style.RenderCheckboxFixedWidth("Restore Previous Settings",
		m.prefsRestoreSettings, m.focus == FocusPrefsRestoreSettings, labelWidth)

	closeStyle := style.NormalButton
	if m.focus == FocusPrefsClose {
		closeStyle = style.FocusButton
	}
	closeBtn := closeStyle.Render(" Close ")

	return lipgloss.JoinVertical(lipgloss.Center,
		desc, "", showHelp, restoreTheme, restoreSettings, "", closeBtn)
}
