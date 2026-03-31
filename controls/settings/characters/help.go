package characters

var helpText = map[State]string{
	Ascii:         "Ascii mode renders with ASCII characters colored by the image.",
	Unicode:       "Block mode uses unicode block characters with foreground and background colors for high detail.",
	Custom:        "Custom mode lets you define your own symbol set for rendering.",
	DarkVariance:  "Maps characters by distance from the darkest color. Darker areas get earlier characters in the set.",
	LightVariance: "Maps characters by distance from the lightest color. Lighter areas get earlier characters in the set.",
	Sequence:      "Cycles through characters sequentially, repeating across the image.",
	Random:        "Assigns characters randomly using a deterministic seed.",
	SeedForm:      "The random seed determines which characters appear where. Same seed = same pattern.",
	ThresholdForm: "Variance below this threshold renders as a space instead of a character.",
}

// HelpText returns the help tooltip for the currently focused item.
func (m Model) HelpText() string {
	if text, ok := helpText[m.focus]; ok {
		return text
	}
	// Check ascii/unicode sub-modes
	if _, ok := asciiCharModeMap[m.focus]; ok {
		return "Select which ASCII characters to use for rendering."
	}
	if _, ok := unicodeCharModeMap[m.focus]; ok {
		return "Select which unicode block character style to use."
	}
	return ""
}
