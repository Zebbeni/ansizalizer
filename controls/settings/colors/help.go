package colors

// HelpText returns contextual help for the currently focused colors item.
func (m Model) HelpText() string {
	if m.focus == Palette {
		return m.PaletteControls.HelpText()
	}
	switch m.focus {
	case UseTrueColor:
		return "True Color uses the full 24-bit RGB spectrum (16 million colors) for maximum fidelity."
	case UsePalette:
		return "Palette mode maps image colors to a limited set. Choose palettes from Lospec, files, or generate from the image."
	case AdaptOn:
		return "Adapt stretches the image's color range to better match the palette. Improves contrast when the palette doesn't cover the image's range well."
	case AdaptOff:
		return "Without adapting, image colors are matched directly to the nearest palette color."
	}
	return ""
}
