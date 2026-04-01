package textstyle

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/style"
)

type State int

const (
	BoldOn State = iota
	BoldOff
	ItalicOn
	ItalicOff
	UnderlineOn
	UnderlineOff
	StrikethroughOn
	StrikethroughOff
)

type Model struct {
	focus         State
	bold          bool
	italic        bool
	underline     bool
	strikethrough bool
	ShouldClose   bool
	IsActive      bool
	width         int
}

func New(w int) Model {
	return Model{
		focus:    BoldOff,
		IsActive: false,
		width:    w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			m.IsActive = false
			m.ShouldClose = true
			return m, nil
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, nil
}

func (m Model) View() string {
	rows := []string{
		m.drawToggle("Bold", BoldOn, BoldOff, m.bold),
		m.drawToggle("Italic", ItalicOn, ItalicOff, m.italic),
		m.drawToggle("Underline", UnderlineOn, UnderlineOff, m.underline),
		m.drawToggle("Strikethrough", StrikethroughOn, StrikethroughOff, m.strikethrough),
	}
	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(content)
}

func (m Model) Bold() bool          { return m.bold }
func (m Model) Italic() bool        { return m.italic }
func (m Model) Underline() bool     { return m.underline }
func (m Model) Strikethrough() bool { return m.strikethrough }

func (m *Model) SetConfig(bold, italic, underline, strikethrough bool) {
	m.bold = bold
	m.italic = italic
	m.underline = underline
	m.strikethrough = strikethrough
}

func (m *Model) ResetFocus() {
	m.focus = BoldOn
}

func (m Model) Summary() string {
	active := make([]string, 0, 4)
	if m.bold {
		active = append(active, "Bold")
	}
	if m.italic {
		active = append(active, "Italic")
	}
	if m.underline {
		active = append(active, "Underline")
	}
	if m.strikethrough {
		active = append(active, "Strikethrough")
	}
	if len(active) == 0 {
		return "None"
	}
	result := ""
	for i, a := range active {
		if i > 0 {
			result += ", "
		}
		result += a
	}
	return result
}

func (m Model) drawToggle(title string, onState, offState State, value bool) string {
	focused := m.IsActive && (m.focus == onState || m.focus == offState)
	return lipgloss.NewStyle().PaddingLeft(1).Render(
		style.RenderCheckboxFixedWidth(title, value, focused, 15))
}

func (m Model) DebugState() string {
	return fmt.Sprintf("TextStyle: bold=%v italic=%v underline=%v strikethrough=%v",
		m.bold, m.italic, m.underline, m.strikethrough)
}
