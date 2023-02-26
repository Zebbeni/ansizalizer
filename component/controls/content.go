package controls

import (
	"github.com/Zebbeni/ansizalizer/component/browser"
	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type addContent func(Content)

type Content interface {
	tea.Model

	// GetActivePosition returns the position of the active or focused content
	// as a percentage of the content's total width and height.
	// This allows a viewport to adjust its Offsets to ensure the active area
	// is still shown if the full rendered string is larger than its max width
	// or height
	// For example, a menu with 100 items and an active item index of 5 might
	// return (0, .05).
	GetActivePosition() (float64, float64)
}

func (m Model) addContent(content Content) {
	m.content = append(m.content, content)
}

func (m Model) removeContent() {
	if len(m.content) > 1 {
		m.content = m.content[:len(m.content)-1]
	}
}

func (m Model) BuildMainMenu() Content {
	updateDir := func(dir string) { m.navState.DirPath = dir }
	updateFile := func(file string) { m.navState.Filepath = file }

	return menu.New([]item.Model{
		item.New("Open", func() {
			f := browser.New(updateDir, updateFile)
			m.addContent(f)
		}),
		item.New("Settings", func() {}),
		item.New("Process", func() {}),
	})
}
