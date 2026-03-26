package adaptive

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder = []State{CountForm, IterForm}
	stateNames = map[State]string{
		CountForm: "Colors",
		IterForm:  "Passes",
	}

	inputStyle = style.BgStyle().Width(13).AlignHorizontal(lipgloss.Left).PaddingLeft(1)

)

func (m Model) drawTitle() string {
	title := style.DimmedTitle.Copy().Italic(true).Render("Generate palette From image")
	return style.BgStyle().Width(m.width).PaddingBottom(1).AlignHorizontal(lipgloss.Center).Render(title)
}

func (m Model) drawInputs() string {
	prompt, placeholder := m.getInputColors(CountForm)

	m.countInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.countInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.countInput.Focused() {
		m.countInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.countInput.Cursor.SetMode(cursor.CursorHide)
	}

	prompt, placeholder = m.getInputColors(IterForm)
	m.iterInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.iterInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.iterInput.Focused() {
		m.iterInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.iterInput.Cursor.SetMode(cursor.CursorHide)
	}

	countInput := inputStyle.Render(m.countInput.View())
	iterInput := inputStyle.Render(m.iterInput.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, countInput, iterInput)
}

func (m Model) drawGenerateButton() string {
	styleColor := style.DimmedColor1
	if m.IsActive && m.focus == Generate {
		styleColor = style.SelectedColor1
	} else if m.active == Generate {
		styleColor = style.NormalColor1
	}

	btnStyle := style.BgStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styleColor).
		Foreground(styleColor)

	button := btnStyle.Render("Generate New")
	return style.BgStyle().Width(m.width).PaddingLeft(1).Render(button)
}

// TODO: This is almost the same as drawGenerateButton. See if we can generalize
func (m Model) drawSaveButton() string {
	styleColor := style.DimmedColor1
	if m.IsActive && m.focus == Save {
		styleColor = style.SelectedColor1
	} else if m.active == Save {
		styleColor = style.NormalColor1
	}

	btnStyle := style.BgStyle().
		Width(m.width - 4).
		AlignHorizontal(lipgloss.Center).
		PaddingTop(1).
		Foreground(styleColor)

	button := btnStyle.Render("Save to .hex File")
	return style.BgStyle().Width(m.width).PaddingLeft(1).Render(button)
}

func (m Model) getInputColors(state State) (lipgloss.TerminalColor, lipgloss.TerminalColor) {
	if m.IsActive {
		if m.focus == state {
			return style.SelectedColor1, style.SelectedColor1
		} else if m.active == state {
			return style.NormalColor1, style.NormalColor1
		}
	}
	return style.DimmedColor1, style.DimmedColor1
}
