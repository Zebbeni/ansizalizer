package controls

import (
	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
	"github.com/Zebbeni/ansizalizer/component/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type Content interface {
	tea.Model
	keyboard.Handler
}

func NewMainMenu(k *keyboard.Map) Content {
	return menu.New([]item.Model{
		item.New("Open", func() {}),
		item.New("Settings", func() {}),
		item.New("Process", func() {}),
	}, k)
}
