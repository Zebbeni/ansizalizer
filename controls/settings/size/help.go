package size

// HelpText returns contextual help for the currently focused size item.
func (m Model) HelpText() string {
	switch m.focus {
	case FitButton:
		return "Fit preserves aspect ratio within the width and height limits."
	case FillButton:
		return "Fill preserves aspect ratio but may crop to fill the entire area."
	case StretchButton:
		return "Stretch distorts the image to exactly fill width and height."
	case WidthForm:
		return "Output width in characters."
	case HeightForm:
		return "Output height in character rows."
	case CharRatioForm:
		return "Width-to-height ratio of your terminal's font. Adjust if output looks stretched. Typical: 0.45-0.50."
	}
	return ""
}
