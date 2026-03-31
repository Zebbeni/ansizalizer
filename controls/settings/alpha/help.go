package alpha

// HelpText returns contextual help for the currently focused alpha item.
func (m Model) HelpText() string {
	switch m.focus {
	case UseAlpha:
		return "Toggle transparency. When enabled, pixels below the alpha threshold render as empty space."
	case TrimAlpha:
		return "Trim removes leading/trailing transparent rows and columns from the output."
	case ThresholdForm:
		return "Pixels with opacity below this threshold are treated as transparent. 0 = only fully transparent pixels, 1 = all pixels treated as transparent."
	}
	return ""
}
