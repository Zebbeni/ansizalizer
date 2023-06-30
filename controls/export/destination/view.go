package destination

import (
	"path/filepath"

	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawSource() string {
	valueStyle := style.DimmedTitle
	if DstInput == m.focus {
		valueStyle = style.SelectedTitle
	}

	value := m.Browser.ActiveFile
	if m.doExportDirectory {
		value = m.SourceBrowser.ActiveDir
	}
	value = filepath.Base(value)

	if value == "." || len(value) == 0 {
		value = "(None)"
	}

	valueContent := valueStyle.Render(value)

	valueWidth := m.width
	widthStyle := lipgloss.NewStyle().Width(valueWidth).AlignHorizontal(lipgloss.Center)
	valueContent = widthStyle.Render(valueContent)

	if m.focus != SrcBrowser {
		return valueContent
	}

	browserContent := m.SourceBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Left, valueContent, browserContent)
}
