package basic

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

const (
	maxWidth          = 23
	maxNormalHeight   = 1
	maxSelectedHeight = 2
)

// NewItemStyles returns style definitions for a default item.
// DefaultItemView for when these come into play.
func NewItemStyles() (s list.DefaultItemStyles) {

	s.NormalTitle = style.NormalTitle.Copy().Padding(0, 0, 0, 2)
	s.NormalDesc = style.NormalParagraph.Copy().MaxHeight(maxNormalHeight).Padding(0, 0, 0, 2)

	s.SelectedTitle = style.SelectedTitle.Copy().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Padding(0, 0, 0, 1)
	s.SelectedDesc = style.SelectedTitle.Copy().MaxHeight(maxSelectedHeight).Padding(0, 0, 0, 2)

	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 0, 0, 2)
	s.DimmedDesc = style.DimmedParagraph.Copy().MaxHeight(maxNormalHeight).Padding(0, 0, 0, 2)

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}
