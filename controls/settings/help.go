package settings

// HelpText returns tooltip text for the currently focused settings item.
func (m Model) HelpText() string {
	// If a sub-panel is active, delegate to it
	switch m.active {
	case Characters:
		return m.Characters.HelpText()
	case Colors:
		return m.Colors.HelpText()
	case Advanced:
		return m.Advanced.HelpText()
	case Size:
		return m.Size.HelpText()
	case Alpha:
		return m.Alpha.HelpText()
	case Adjust:
		return m.Adjust.HelpText()
	case TextStyle:
		return m.TextStyle.HelpText()
	case Animation:
		return m.Animation.HelpText()
	}

	// Top-level settings focus
	switch m.focus {
	case Colors:
		return "Choose between true color and palette-based rendering. Select palettes from Lospec, files, or generated from the current image."
	case Characters:
		return "Select character mode: Block (unicode), Ascii, or Custom symbols. Each mode offers different visual styles for rendering."
	case Size:
		return "Set the output width and height in characters, and the character aspect ratio for your terminal font."
	case Adjust:
		return "Fine-tune brightness and contrast of the rendered output."
	case TextStyle:
		return "Apply ANSI text attributes: bold, italic, underline, and strikethrough."
	case Advanced:
		return "Configure the image sampling method and dithering algorithm for palette-based rendering."
	case Alpha:
		return "Control how transparent pixels are rendered. Set a threshold to treat semi-transparent pixels as fully transparent."
	case Animation:
		return "Set the frame delay for animated GIF playback."
	}
	return ""
}
