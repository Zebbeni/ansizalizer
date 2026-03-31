package textstyle

// HelpText returns contextual help for the currently focused text style item.
func (m Model) HelpText() string {
	switch m.focus {
	case BoldOn, BoldOff:
		return "Bold makes characters thicker. Always applied in Ascii and Custom modes."
	case ItalicOn, ItalicOff:
		return "Italic slants characters. Support varies by terminal."
	case UnderlineOn, UnderlineOff:
		return "Underline adds a line beneath characters. Support varies by terminal."
	case StrikethroughOn, StrikethroughOff:
		return "Strikethrough draws a line through characters. Support varies by terminal."
	}
	return ""
}
