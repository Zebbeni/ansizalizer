package browser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/event"
)

type Model struct {
	SelectedDir  string
	SelectedFile string
	ActiveDir    string
	ActiveFile   string

	lists          []list.Model
	fileExtensions map[string]bool

	ShouldClose bool

	width int
}

func New(exts map[string]bool, w int) Model {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting starting directory:", err)
		os.Exit(1)
	}

	m := Model{
		width:          w,
		fileExtensions: exts,
	}
	m = m.addListForDirectory(dir)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Esc):
			return m.handleEsc()
		case key.Matches(msg, event.KeyMap.Nav):
			return m.handleNav(msg)
		case key.Matches(msg, event.KeyMap.Enter):
			return m.handleEnter()
		}
	}
	return m, nil
}

func (m Model) currentList() list.Model {
	return m.lists[m.listIndex()]
}

func (m Model) listIndex() int {
	return len(m.lists) - 1
}

func (m Model) View() string {
	browser := m.currentList().View()
	return lipgloss.JoinVertical(lipgloss.Left, browser)
}

func (m Model) ActiveFilename() string {
	return filepath.Base(m.ActiveFile)
}
