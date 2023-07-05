package destination

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawSelected() string {
	title := style.DimmedTitle.Copy().Render("Selected")

	valueStyle := style.DimmedTitle.Copy()

	if Input == m.focus {
		if m.IsActive {
			valueStyle = style.SelectedTitle.Copy()
		} else {
			valueStyle = style.NormalTitle.Copy()
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
	widthStyle := lipgloss.NewStyle().Width(valueWidth).AlignHorizontal(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, title, valueContent)

	return widthStyle.Render(content)
}

func drawBrowserTitle() string {
	return style.DimmedTitle.Copy().Padding(0, 2, 1, 2).Render("Select a directory")
}
