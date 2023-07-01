package destination

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawInput() string {
	valueStyle := style.DimmedTitle
	if Input == m.focus {
		valueStyle = style.SelectedTitle
	}

	dir := m.Browser.ActiveDir
	parentDir := filepath.Base(filepath.Dir(dir))
	activeDir := filepath.Base(dir)
	value := fmt.Sprintf("%s/%s", parentDir, activeDir)

	if value == "." || len(value) == 0 {
		value = "(None)"
	}

	valueContent := valueStyle.Render(value)

	valueWidth := m.width
	widthStyle := lipgloss.NewStyle().Width(valueWidth).AlignHorizontal(lipgloss.Center)
	valueContent = widthStyle.Render(valueContent)

	if m.focus != Browser {
		return valueContent
	}

	browserContent := m.Browser.View()
	return lipgloss.JoinVertical(lipgloss.Left, valueContent, browserContent)
}
