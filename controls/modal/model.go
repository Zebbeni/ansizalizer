package modal

import (
	"path/filepath"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

const (
	Width  = 50
	Height = 20
)

type Kind int

const (
	LoadKind Kind = iota
	SaveKind
)

type Focus int

const (
	FocusBrowser Focus = iota
	FocusInput
	FocusConfirm
	FocusCancel
)

// Result is set when the user confirms an action.
type Result struct {
	Kind Kind
	Path string // full file path for load or save
}

type Model struct {
	Kind    Kind
	Open    bool
	focus   Focus
	browser browser.Model
	input   textinput.Model
	result  *Result

	dir   string // current directory for save
	width int
}

var jsonExtensions = map[string]bool{".json": true}

func NewLoad(settingsDir string) Model {
	b := browser.NewAtDir(settingsDir, jsonExtensions, Width-4)
	return Model{
		Kind:    LoadKind,
		Open:    true,
		focus:   FocusBrowser,
		browser: b,
		width:   Width,
	}
}

func NewSave(settingsDir string) Model {
	input := textinput.New()
	input.Placeholder = "settings"
	input.Prompt = "Filename: "
	input.Focus()
	input.CharLimit = 30
	input.SetWidth(Width - 10)
	return Model{
		Kind:  SaveKind,
		Open:  true,
		focus: FocusInput,
		input: input,
		dir:   settingsDir,
		width: Width,
	}
}

func (m Model) GetResult() *Result {
	return m.result
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.Open {
		return m, nil
	}

	switch m.Kind {
	case LoadKind:
		return m.updateLoad(msg)
	case SaveKind:
		return m.updateSave(msg)
	}
	return m, nil
}

func (m Model) updateLoad(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == FocusBrowser {
		var cmd tea.Cmd
		m.browser, cmd = m.browser.Update(msg)

		if m.browser.ShouldClose {
			m.browser.ShouldClose = false
			if m.browser.SelectedFile != "" {
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
		if m.focus == FocusConfirm && m.browser.SelectedFile != "" {
			m.result = &Result{Kind: LoadKind, Path: m.browser.SelectedFile}
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
			m.focus = FocusBrowser
			m.browser.SetActive(true)
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.Open = false
	}
	return m, nil
}

func (m Model) updateSave(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == FocusInput {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, event.KeyMap.Enter):
				m.input.Blur()
				m.focus = FocusConfirm
				return m, nil
			case key.Matches(keyMsg, event.KeyMap.Esc):
				m.Open = false
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		if m.focus == FocusConfirm {
			filename := m.input.Value()
			if filename == "" {
				filename = "settings"
			}
			m.result = &Result{Kind: SaveKind, Path: filepath.Join(m.dir, filename+".json")}
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
			m.focus = FocusInput
			m.input.Focus()
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.Open = false
	}
	return m, nil
}

func (m Model) View() string {
	var title, body string

	switch m.Kind {
	case LoadKind:
		title = "Load Settings"
		body = m.loadView()
	case SaveKind:
		title = "Save Settings"
		body = m.saveView()
	}

	titleRendered := style.SelectedTitle.Copy().
		Width(Width - 6).
		AlignHorizontal(lipgloss.Center).
		Render(title)

	buttons := m.renderButtons()

	content := lipgloss.JoinVertical(lipgloss.Center,
		titleRendered,
		"",
		body,
		"",
		buttons,
	)

	border := style.BgStyle().
		Width(Width).
		Height(Height).
		Border(style.HeavyBorder()).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg).
		Padding(1, 2).
		AlignHorizontal(lipgloss.Center)

	return style.ApplyBg(border.Render(content), 0)
}

func (m Model) loadView() string {
	if m.browser.SelectedFile != "" && m.focus != FocusBrowser {
		filename := filepath.Base(m.browser.SelectedFile)
		return style.NormalTitle.Copy().
			Width(Width - 6).
			AlignHorizontal(lipgloss.Center).
			Render("Load " + filename + "?")
	}
	dirLabel := style.DimmedTitle.Copy().Italic(true).
		Width(Width - 6).
		AlignHorizontal(lipgloss.Center).
		Render(filepath.Base(m.browser.SelectedDir) + "/")
	return lipgloss.JoinVertical(lipgloss.Left, dirLabel, m.browser.View())
}

func (m Model) saveView() string {
	promptStyle := style.DimmedTitle.Copy()
	textSt := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	if m.input.Focused() {
		promptStyle = style.SelectedTitle.Copy()
		textSt = lipgloss.NewStyle().Foreground(style.SelectedColor1)
	}

	styles := m.input.Styles()
	styles.Focused.Prompt = promptStyle.Width(12)
	styles.Focused.Text = textSt
	styles.Blurred.Prompt = promptStyle.Width(12)
	styles.Blurred.Text = textSt
	styles.Cursor.Blink = m.input.Focused()
	m.input.SetStyles(styles)

	return m.input.View()
}

func (m Model) renderButtons() string {
	confirmLabel := "  Load  "
	if m.Kind == SaveKind {
		confirmLabel = "  Save  "
	}

	confirmStyle := style.NormalButton
	if m.focus == FocusConfirm {
		confirmStyle = style.FocusButton
	}
	cancelStyle := style.NormalButton
	if m.focus == FocusCancel {
		cancelStyle = style.FocusButton
	}

	return lipgloss.JoinHorizontal(lipgloss.Center,
		confirmStyle.Render(confirmLabel),
		"  ",
		cancelStyle.Render(" Cancel "),
	)
}
