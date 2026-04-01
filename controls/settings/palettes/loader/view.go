package loader

import (
	"path/filepath"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

const (
	maxWidth          = 30
	maxNormalHeight   = 1
	maxSelectedHeight = 2
)

// NewItemStyles returns style definitions for a default item.
// DefaultItemView for when these come into play.
func NewItemStyles() (s list.DefaultItemStyles) {

	s.NormalTitle = style.DimmedTitle.Padding(0, 1, 0, 2)
	s.NormalDesc = style.DimmedParagraph.MaxHeight(maxNormalHeight).Padding(0, 0, 0, 2)

	s.SelectedTitle = style.SelectedTitle.Padding(0, 1, 0, 1).
		Border(style.HeavyBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)
	s.SelectedDesc = style.SelectedTitle.MaxHeight(maxSelectedHeight).Padding(0, 0, 0, 1).
		Border(style.HeavyBorder(), false, false, false, true).
		BorderForeground(style.SelectedColor1)

	s.DimmedTitle = style.DimmedTitle.Padding(0, 1, 0, 0)
	s.DimmedDesc = style.DimmedParagraph.MaxHeight(maxNormalHeight).Padding(0, 0, 0, 2)

	return s
}

func (m Model) drawTitle() string {
	dir := filepath.Base(m.FileBrowser.SelectedDir)
	title := style.DimmedTitle.Italic(true).Render("Browsing " + dir + "/")
	return style.BgStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}
