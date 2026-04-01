package destination

import (
	"fmt"
	"path/filepath"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawSelected() string {
	title := style.DimmedTitle.Render("Selected")

	valueStyle := style.DimmedTitle

	if Input == m.focus {
		if m.IsActive {
			valueStyle = style.SelectedTitle
		} else {
			valueStyle = style.NormalTitle
		}
	}
	valueStyle.Padding(0, 0, 1, 0)

	path := m.Browser.SelectedDir

	parent := filepath.Base(filepath.Dir(path))
	selected := filepath.Base(path)
	value := fmt.Sprintf("%s/%s", parent, selected)

	valueRunes := []rune(value)
	if len(valueRunes) > m.width {
		value = string(valueRunes[len(valueRunes)-m.width:])
	}

	valueContent := valueStyle.Render(value)

	valueWidth := m.width
	widthStyle := style.BgStyle().Width(valueWidth).AlignHorizontal(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, title, valueContent)

	return widthStyle.Render(content)
}

func (m Model) drawBrowserTitle() string {
	dir := filepath.Base(m.Browser.SelectedDir)
	return style.DimmedTitle.Italic(true).Padding(0, 2, 1, 2).Render("Browsing " + dir + "/")
}
