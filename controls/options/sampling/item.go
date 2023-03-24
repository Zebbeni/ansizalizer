package sampling

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/nfnt/resize"
)

type item struct {
	name     string
	Function resize.InterpolationFunction
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

func menuItems() []list.Item {
	items := make([]list.Item, len(nameMap))
	for i, f := range Functions {
		items[i] = item{name: nameMap[f], Function: f}
	}
	return items
}

func newMenu(items []list.Item, height int) list.Model {
	l := list.New(items, NewDelegate(), 30, height)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)

	l.KeyMap.ForceQuit.Unbind()
	l.KeyMap.Quit.Unbind()
	return l
}

func NewDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	return delegate
}
