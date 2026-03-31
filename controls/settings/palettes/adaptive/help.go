package adaptive

// HelpText returns contextual help for the currently focused adaptive palette item.
func (m Model) HelpText() string {
	switch m.focus {
	case CountForm:
		return "Number of colors to extract from the image."
	case IterForm:
		return "More passes = more accurate palette, but slower."
	case Generate:
		return "Generate a new palette by sampling the current image's colors."
	case Save:
		return "Save the generated palette to a .hex file."
	}
	return ""
}
