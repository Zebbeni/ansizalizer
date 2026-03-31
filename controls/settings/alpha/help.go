package alpha

// HelpText returns contextual help for the currently focused alpha item.
func (m Model) HelpText() string {
	switch m.focus {
	case AlphaYes:
		return "Enable alpha channel support. Transparent pixels render as empty space."
	case AlphaNo:
		return "Disable alpha. All pixels render as opaque, ignoring transparency."
	case TrimAlphaYes, TrimAlphaNo:
		return "Trim removes leading/trailing transparent rows and columns from the output."
	case ThresholdForm:
		return "Pixels with opacity below this threshold are treated as transparent. 0 = only fully transparent, 0.5 = half-transparent and below."
	}
	return ""
}
