package style

import "image/color"

// InputColors returns prompt and text colors for a numeric/text input field
// based on whether the component is active, focused, and editing.
// This replaces the duplicated getInputColors pattern across settings packages.
func InputColors(isActive, isFocused, isEditing bool) (color.Color, color.Color) {
	if isActive && isFocused {
		if isEditing {
			return NormalColor1, SelectedColor1
		}
		return SelectedColor1, NormalColor1
	}
	return DimmedColor1, DimmedColor1
}
