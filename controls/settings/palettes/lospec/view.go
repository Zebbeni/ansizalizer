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
	bg := style.BgStyle()

	m.countInput.SetWidth(3)
	cStyles := m.countInput.Styles()
	cStyles.Focused.Prompt = bg.Copy().Foreground(prompt)
	cStyles.Focused.Text = bg.Copy().Foreground(prompt).MaxWidth(3)
	cStyles.Focused.Placeholder = bg.Copy().Foreground(placeholder)
	cStyles.Blurred.Prompt = bg.Copy().Foreground(prompt)
	cStyles.Blurred.Text = bg.Copy().Foreground(prompt).MaxWidth(3)
	cStyles.Blurred.Placeholder = bg.Copy().Foreground(placeholder)
	cStyles.Cursor.Blink = m.countInput.Focused()
	cStyles.Cursor.Color = prompt
	m.countInput.SetStyles(cStyles)
	m.countInput.SetVirtualCursor(m.countInput.Focused())
	return bg.Copy().PaddingLeft(1).Width(12).Render(m.countInput.View())
}

func (m Model) drawTagInput() string {
	prompt, placeholder := m.getInputColors(TagForm)
	bg := style.BgStyle()

	m.tagInput.SetWidth(m.width - 5)
	tStyles := m.tagInput.Styles()
	tStyles.Focused.Prompt = bg.Copy().Foreground(prompt).Padding(0, 1, 0, 1)
	tStyles.Focused.Text = bg.Copy().Foreground(prompt)
	tStyles.Focused.Placeholder = bg.Copy().Foreground(placeholder)
	tStyles.Blurred.Prompt = bg.Copy().Foreground(prompt).Padding(0, 1, 0, 1)
	tStyles.Blurred.Text = bg.Copy().Foreground(prompt)
	tStyles.Blurred.Placeholder = bg.Copy().Foreground(placeholder)
	tStyles.Cursor.Blink = m.tagInput.Focused()
	tStyles.Cursor.Color = prompt
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
