package menu

import (
	"github.com/charmbracelet/bubbles/list"
)

func New(items []list.Item) list.Model {
	newList := list.New(items, NewDelegate(), 25, 20)

	newList.KeyMap.ForceQuit.Unbind()
	newList.KeyMap.Quit.Unbind()
	newList.SetShowHelp(false)
	newList.SetShowStatusBar(false)
	newList.SetShowTitle(false)
	newList.SetFilteringEnabled(false)

	return newList
}

func NewDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	return delegate
}
