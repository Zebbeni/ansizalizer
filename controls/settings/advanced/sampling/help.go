package sampling

// HelpText returns contextual help for the sampling list.
func (m Model) HelpText() string {
	return "Select the resize interpolation method. Nearest Neighbor is fastest; Lanczos3 is highest quality."
}
