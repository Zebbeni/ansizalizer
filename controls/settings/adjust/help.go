package adjust

// HelpText returns contextual help for the currently focused adjust item.
func (m Model) HelpText() string {
	switch m.focus {
	case BrightnessForm:
		return "Brightness offset (-100 to 100). Positive = lighter, negative = darker."
	case ContrastForm:
		return "Contrast adjustment (-100 to 100). Higher values increase the difference between light and dark."
	}
	return ""
}
