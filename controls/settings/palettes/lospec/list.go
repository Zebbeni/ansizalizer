package lospec

import (
	"charm.land/bubbles/v2/list"

	"github.com/Zebbeni/ansizalizer/style"
)

func CreateList(items []list.Item, w int) list.Model {
	newList := list.New(items, NewDelegate(), w, 22)

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
	delegate.ShowDescription = true
	delegate.Styles = ItemStyles()
	return delegate
}

func ItemStyles() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Padding(0, 1, 0, 2)
	s.NormalDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	s.SelectedTitle = style.SelectedTitle.Padding(0, 1, 0, 1).
		Border(style.HeavyBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)
	s.SelectedDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	s.DimmedTitle = style.DimmedTitle.Padding(0, 1, 0, 0)
	s.DimmedDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	return s
}

func InactiveItemStyles() (s list.DefaultItemStyles) {
	s.NormalTitle = style.DimmedTitle.Padding(0, 1, 0, 2)
	s.NormalDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	s.SelectedTitle = style.NormalTitle.Padding(0, 1, 0, 2)
	s.SelectedDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	s.DimmedTitle = style.DimmedTitle.Padding(0, 1, 0, 0)
	s.DimmedDesc = style.DimmedParagraph.MaxHeight(1).Padding(0, 0, 0, 2)

	return s
}
