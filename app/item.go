package app

import "github.com/charmbracelet/bubbles/list"

type item struct {
	name  string
	state State
}

func (i item) FilterValue() string {
	return i.name
}

func (i item) Title() string {
	return i.name
}

func (i item) Description() string {
	return ""
}

func newMenu() list.Model {
	items := []list.Item{
		item{name: "File", state: Browser},
		item{name: "Options", state: Settings},
	}
	menu := list.New(items, NewDelegate(), 20, 20)
	menu.SetShowHelp(false)
	menu.SetShowFilter(false)
	menu.SetShowTitle(false)
	menu.SetShowStatusBar(false)

	menu.KeyMap.ForceQuit.Unbind()
	menu.KeyMap.Quit.Unbind()
	return menu
}

func NewDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	return delegate
}
