package basic

import (
	"github.com/charmbracelet/bubbles/list"

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

	s.NormalTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 1)
	s.NormalDesc = style.DimmedParagraph.Copy().MaxHeight(maxNormalHeight).Padding(0, 1, 0, 1)

	s.SelectedTitle = style.SelectedTitle.Copy().Padding(0, 1, 0, 1)
	//Border(lipgloss.NormalBorder(), false, false, false, true).
	//BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
	s.SelectedDesc = style.SelectedTitle.Copy().MaxHeight(maxSelectedHeight).Padding(0, 1, 0, 1)

	s.DimmedTitle = style.DimmedTitle.Copy().Padding(0, 1, 0, 0)
	s.DimmedDesc = style.DimmedParagraph.Copy().MaxHeight(maxNormalHeight).Padding(0, 1, 0, 1)

	//s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}
