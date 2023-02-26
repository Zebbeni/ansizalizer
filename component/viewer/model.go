package viewer

import (
	"github.com/Zebbeni/ansizalizer/component/style"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// rendering
	w, h int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m Model) View() string {
	imageString := "Viewer" // eventually we'll display ansi stuff here

	content := lipgloss.NewStyle().Width(m.w).Height(m.h).Render(imageString)

	vp := viewport.New(m.w, m.h)
	vp.SetContent(content)
	vp.Style = style.ViewerBorder.Copy().Width(m.w).Height(m.h)

	return vp.View()
}

func (m Model) Resize(w, h int) {
	m.w, m.h = w, h
}
