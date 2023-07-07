package sampling

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/nfnt/resize"

	"github.com/Zebbeni/ansizalizer/style"
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

func newMenu(items []list.Item, width, height int) list.Model {
	l := list.New(items, NewDelegate(false), width, height)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)

	l.KeyMap.ForceQuit.Unbind()
	l.KeyMap.Quit.Unbind()

	return l
}

func NewDelegate(isActive bool) list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	if isActive {
		delegate.Styles = ItemStylesActive()
	} else {
		delegate.Styles = ItemStylesInactive()
	}
	return delegate
}

func ItemStylesActive() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 2)
	s.SelectedTitle = style.SelectedTitle.Copy().Padding(0, 1, 0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)
	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	return s
}

func ItemStylesInactive() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 2)
	s.SelectedTitle = style.NormalTitle.Copy().Padding(0, 1, 0, 2)
	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	return s
}
