package menu

import (
	"charm.land/bubbles/v2/list"

	"github.com/Zebbeni/ansizalizer/style"
)

func New(items []list.Item, w int) list.Model {
	newList := list.New(items, NewDelegate(), w, 18)

	newList.KeyMap.ForceQuit.Unbind()
	newList.KeyMap.Quit.Unbind()
	newList.SetShowHelp(false)
	newList.SetShowStatusBar(false)
	newList.SetShowTitle(false)
	newList.SetFilteringEnabled(false)
	newList.Styles.NoItems = style.BgStyle()
	newList.Styles.PaginationStyle = style.BgStyle()

	return newList
}

func NewDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	delegate.Styles = ItemStyles()
	return delegate
}

func ItemStyles() (s list.DefaultItemStyles) {
	bg := style.BgStyle()

	s.NormalTitle = bg.Foreground(style.DimmedColor1).Padding(0, 1, 0, 2)
	s.NormalDesc = bg.Foreground(style.DimmedColor2).MaxHeight(1).Padding(0, 0, 0, 2)

	s.SelectedTitle = bg.Foreground(style.SelectedColor1).Padding(0, 1, 0, 1).
		Border(style.HeavyBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg)
	s.SelectedDesc = bg.Foreground(style.SelectedColor1).MaxHeight(1).Padding(0, 0, 0, 1).
		Border(style.HeavyBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg)

	s.DimmedTitle = bg.Foreground(style.DimmedColor1).Padding(0, 1, 0, 0)
	s.DimmedDesc = bg.Foreground(style.DimmedColor2).MaxHeight(1).Padding(0, 0, 0, 2)

	return s
}

func InactiveItemStyles() (s list.DefaultItemStyles) {
	bg := style.BgStyle()

	s.NormalTitle = bg.Foreground(style.DimmedColor1).Padding(0, 1, 0, 2)
	s.NormalDesc = bg.Foreground(style.DimmedColor2).MaxHeight(1).Padding(0, 0, 0, 2)

	s.SelectedTitle = bg.Foreground(style.NormalColor1).Padding(0, 1, 0, 2)
	s.SelectedDesc = bg.Foreground(style.DimmedColor2).MaxHeight(1).Padding(0, 0, 0, 2)

	s.DimmedTitle = bg.Foreground(style.DimmedColor1).Padding(0, 1, 0, 0)
	s.DimmedDesc = bg.Foreground(style.DimmedColor2).MaxHeight(1).Padding(0, 0, 0, 2)

	return s
}
