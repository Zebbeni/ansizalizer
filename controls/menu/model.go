package menu

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

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
	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 2)
	s.NormalDesc = style.DimmedParagraph.Copy().MaxHeight(1).Padding(0, 0, 0, 2)

	s.SelectedTitle = style.SelectedTitle.Copy().Padding(0, 1, 0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)
	s.NormalDesc = style.SelectedTitle.Copy().MaxHeight(1).Padding(0, 0, 0, 2).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)

	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	s.DimmedDesc = style.DimmedParagraph.Copy().MaxHeight(1).Padding(0, 0, 0, 2)

	return s
}
