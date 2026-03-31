package lospec

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateNames = map[State]string{
		CountForm:        "Count",
		TagForm:          "Tag",
		FilterExact:      "Exact",
		FilterMax:        "Max",
		FilterMin:        "Min",
		SortAlphabetical: "A-Z",
		SortDownloads:    "Downloads",
		SortNewest:       "Newest",
	}

	filterOrder = []State{FilterExact, FilterMax, FilterMin}
	sortOrder   = []State{SortDownloads, SortAlphabetical, SortNewest}

)

func (m Model) drawInputs() string {
	colorsInput := m.drawColorsInput()
	tagInput := m.drawTagInput()

	return lipgloss.JoinHorizontal(lipgloss.Left, colorsInput, tagInput)
}

func (m Model) drawTitle() string {
	title := style.DimmedTitle.Copy().Italic(true).Render("Browse Lospec.com")
	return style.BgStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}

func (m Model) drawColorsInput() string {
	prompt, placeholder := m.getInputColors(CountForm)

	m.countInput.CharLimit = 3
	m.countInput.SetWidth(3)
	cStyles := m.countInput.Styles()
	cStyles.Focused.Prompt = cStyles.Focused.Prompt.Foreground(prompt)
	cStyles.Focused.Text = cStyles.Focused.Text.Foreground(prompt).MaxWidth(3)
	cStyles.Focused.Placeholder = cStyles.Focused.Placeholder.Foreground(placeholder)
	cStyles.Blurred.Prompt = cStyles.Blurred.Prompt.Foreground(prompt)
	cStyles.Blurred.Text = cStyles.Blurred.Text.Foreground(prompt).MaxWidth(3)
	cStyles.Blurred.Placeholder = cStyles.Blurred.Placeholder.Foreground(placeholder)
	cStyles.Cursor.Blink = m.countInput.Focused()
	m.countInput.SetStyles(cStyles)
	m.countInput.SetVirtualCursor(m.countInput.Focused())
	return style.BgStyle().Width(11).Render(m.countInput.View())
}

func (m Model) drawTagInput() string {
	prompt, placeholder := m.getInputColors(TagForm)

	m.tagInput.SetWidth(m.width - 5)
	tStyles := m.tagInput.Styles()
	tStyles.Focused.Prompt = tStyles.Focused.Prompt.Foreground(prompt)
	tStyles.Focused.Placeholder = tStyles.Focused.Placeholder.Foreground(placeholder)
	tStyles.Blurred.Prompt = tStyles.Blurred.Prompt.Foreground(prompt)
	tStyles.Blurred.Placeholder = tStyles.Blurred.Placeholder.Foreground(placeholder)
	tStyles.Cursor.Blink = m.tagInput.Focused()
	m.tagInput.SetStyles(tStyles)
	m.tagInput.SetVirtualCursor(m.tagInput.Focused())
	return m.tagInput.View()
}

func (m Model) drawFilterButtons() string {
	buttons := make([]string, len(filterOrder))
	for i, filter := range filterOrder {
		buttonStyle := style.NormalButtonNode
		if filter == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if filter == m.filterType {
			buttonStyle = style.ActiveButtonNode
		}
		buttons[i] = buttonStyle.Render(stateNames[filter])
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m Model) drawSortButtons() string {
	buttons := make([]string, len(sortOrder))
	for i, sort := range sortOrder {
		buttonStyle := style.NormalButtonNode
		if sort == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if sort == m.sortType {
			buttonStyle = style.ActiveButtonNode
		}
		name := stateNames[sort]
		if sort == m.sortType {
			name = "↑" + name
		} else {
			name = name + " "
		}
		buttons[i] = buttonStyle.Render(name)
	}
	buttonContent := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	return lipgloss.JoinHorizontal(lipgloss.Left, buttonContent)
}

func (m Model) drawPaletteList() string {
	if len(m.paletteList.Items()) == 0 {
		return ""
	}

	return m.paletteList.View()
}

func (m Model) getInputColors(state State) (color.Color, color.Color) {
	if m.IsActive {
		if m.focus == state {
			return style.SelectedColor1, style.SelectedColor1
		} else if m.active == state {
			return style.NormalColor1, style.NormalColor1
		}
	}
	return style.DimmedColor1, style.DimmedColor1
}
