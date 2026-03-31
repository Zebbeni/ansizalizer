package palettes

// HelpText returns contextual help for the currently focused palette item.
func (m Model) HelpText() string {
	switch m.focus {
	case Lospec:
		return "Search Lospec.com for user-created color palettes, filtered by color count and tags."
	case Load:
		return "Load a palette from a local .hex file. Each line is one hex color."
	case Adapt:
		return "Generate a palette by sampling the most prominent colors from the current image."
	}

	// Delegate to active sub-panel
	switch m.controls {
	case Lospec:
		return m.Lospec.HelpText()
	case Adapt:
		return m.Adapter.HelpText()
	}

	return ""
}
