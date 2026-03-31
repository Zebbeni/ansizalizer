package dithering

// HelpText returns contextual help for the currently focused dithering item.
func (m Model) HelpText() string {
	switch m.focus {
	case DitherOn, DitherOff:
		return "Enable or disable dithering. Dithering reduces color banding in palette mode."
	case SerpentineOn, SerpentineOff:
		return "Serpentine alternates error diffusion direction each row, reducing directional artifacts."
	case StrengthForm:
		return "Dither strength controls how aggressively the dithering pattern is applied. 1.0 is default."
	case ModeMatrix:
		return "Error diffusion matrices spread quantization error to neighboring pixels."
	case ModeBayer:
		return "Bayer dithering uses an ordered threshold matrix. Good for a retro, patterned look."
	case ModeClusteredDot:
		return "Clustered dot dithering simulates halftone printing patterns."
	case MatrixList:
		return "Choose an error diffusion algorithm. Floyd-Steinberg is the most common."
	case BayerSize2, BayerSize4, BayerSize8, BayerSize16:
		return "Bayer matrix size. Larger = more gradual transitions, smaller = more visible pattern."
	case ClusteredDotList:
		return "Choose a clustered dot pattern for halftone-style dithering."
	}
	return ""
}
