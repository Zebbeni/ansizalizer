package controls

import (
	"github.com/Zebbeni/ansizalizer/component/browser"
	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
	"github.com/Zebbeni/ansizalizer/component/menu"
	"github.com/Zebbeni/ansizalizer/state"
	tea "github.com/charmbracelet/bubbletea"
)

type updateContent func(Content)
type updateBrowser func(browser *state.Browser)

type Content interface {
	tea.Model
	keyboard.Handler

	// GetActivePosition returns the position of the active or focused content
	// as a percentage of the content's total width and height.
	// This allows a viewport to adjust its Offsets to ensure the active area
	// is still shown if the full rendered string is larger than its max width
	// or height
	// For example, a menu with 100 items and an active item index of 5 might
	// return (0, .05).
	GetActivePosition() (float64, float64)
}

func NewMainMenu(s *state.Model, k *keyboard.Map, update updateContent) Content {
	return menu.New([]item.Model{
		item.New("Open", func() {
			f := browser.New(s.Browser, k)
			update(f)
		}),
		item.New("Settings", func() {}),
		item.New("Process", func() {}),
	}, k)
}
