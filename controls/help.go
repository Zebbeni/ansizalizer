package controls

// HelpText returns the tooltip help text for the currently focused item,
// or empty string if no help is defined.
func (m Model) HelpText() string {
	if m.active == Settings {
		return m.Settings.HelpText()
	}
	switch m.focus {
	case FileMenu:
		return "Open the file menu to load/save settings, change themes, export, or quit."
	case Browse:
		return "Browse for image files to render as ANSI art."
	case Settings:
		return "Adjust rendering settings."
	}
	return ""
}
