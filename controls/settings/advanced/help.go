package advanced

// HelpText returns contextual help for the currently focused advanced item.
func (m Model) HelpText() string {
	if m.active == SamplingControls {
		return m.sampling.HelpText()
	}
	if m.active == DitheringControls {
		return m.dithering.HelpText()
	}
	switch m.focus {
	case Sampling:
		return "The interpolation method used when resizing the image to the target character dimensions."
	case Dithering:
		return "Dithering adds noise patterns to reduce color banding when using limited palettes."
	}
	return ""
}
