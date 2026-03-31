package colors

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawToggles() string {
	paletteFocused := m.IsActive && (m.focus == UseTrueColor || m.focus == UsePalette)
	paletteBox := style.RenderCheckbox("Limited Palette", m.mode == UsePalette, paletteFocused)

	if m.mode != UsePalette {
		return lipgloss.NewStyle().PaddingLeft(1).Render(paletteBox)
	}

	adaptFocused := m.IsActive && (m.focus == AdaptOn || m.focus == AdaptOff)
	adaptBox := style.RenderCheckbox("Adapt", m.adaptToPalette, adaptFocused)

	return lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Left, paletteBox, "  ", adaptBox))
}
