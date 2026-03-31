package colors

// HelpText returns contextual help for the currently focused colors item.
func (m Model) HelpText() string {
	if m.focus == Palette {
		return m.PaletteControls.HelpText()
	}
	switch m.focus {
	case UseTrueColor, UsePalette:
		return "Toggle between full 24-bit color and a limited palette. Palettes can be loaded from Lospec, files, or generated from the image."
	case AdaptOn, AdaptOff:
		return "Adapt stretches the image's color range to better match the palette. Improves contrast when the palette doesn't cover the image's range well."
	}
	return ""
}
