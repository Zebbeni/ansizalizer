package lospec

// HelpText returns contextual help for the currently focused Lospec item.
func (m Model) HelpText() string {
	switch m.focus {
	case CountForm:
		return "Filter palettes by number of colors."
	case TagForm:
		return "Filter by tags (e.g. 'retro', 'pastel', 'gameboy')."
	case FilterExact, FilterMax, FilterMin:
		return "Choose whether the color count is exact, a maximum, or a minimum."
	case SortAlphabetical, SortDownloads, SortNewest:
		return "Sort search results by name, popularity, or date."
	case List:
		return "Browse search results. Press Enter to select a palette."
	}
	return ""
}
