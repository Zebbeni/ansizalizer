package style

import "charm.land/lipgloss/v2"

// RenderCheckbox renders a label followed by a toggleable checkbox.
// checked = current state, focused = whether this item has keyboard focus.
func RenderCheckbox(label string, checked, focused bool) string {
	checkChar := "☐"
	if checked {
		checkChar = "🗹"
	}
	labelStyle := DimmedTitle.Copy()
	checkStyle := DimmedTitle.Copy()
	if focused {
		labelStyle = NormalTitle.Copy()
		checkStyle = SelectedTitle.Copy()
	}
	return lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render(label+" "),
		checkStyle.Render(checkChar))
}

// RenderCheckboxFixedWidth renders a checkbox with a fixed-width label for alignment.
func RenderCheckboxFixedWidth(label string, checked, focused bool, labelWidth int) string {
	checkChar := "☐"
	if checked {
		checkChar = "🗹"
	}
	labelStyle := DimmedTitle.Copy().Width(labelWidth)
	checkStyle := DimmedTitle.Copy()
	if focused {
		labelStyle = NormalTitle.Copy().Width(labelWidth)
		checkStyle = SelectedTitle.Copy()
	}
	return lipgloss.JoinHorizontal(lipgloss.Left,
		labelStyle.Render(label),
		checkStyle.Render(checkChar))
}
