package colors

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawPaletteToggles() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Mode:")

	trueColorStyle := style.NormalButtonNode
	if m.IsActive && m.focus == UseTrueColor {
		trueColorStyle = style.FocusButtonNode
	} else if m.mode == UseTrueColor {
		trueColorStyle = style.ActiveButtonNode
	}
	trueColorNode := trueColorStyle.Render("True Color")
	trueColorNode = lipgloss.NewStyle().PaddingLeft(1).Render(trueColorNode)

	palettedStyle := style.NormalButtonNode
	if m.IsActive && m.focus == UsePalette {
		palettedStyle = style.FocusButtonNode
	} else if m.mode == UsePalette {
		palettedStyle = style.ActiveButtonNode
	}
	palettedNode := palettedStyle.Render("Palette")
	palettedNode = lipgloss.NewStyle().PaddingLeft(1).Render(palettedNode)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, trueColorNode, palettedNode)
}
