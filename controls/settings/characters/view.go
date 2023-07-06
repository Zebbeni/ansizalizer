package characters

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder         = []State{Ascii, Unicode, Custom}
	asciiButtonOrder   = []State{AsciiAz, AsciiNums, AsciiSpec, AsciiAll}
	unicodeButtonOrder = []State{UnicodeFull, UnicodeHalf, UnicodeQuart, UnicodeShadeLight, UnicodeShadeMed, UnicodeShadeHeavy}

	stateNames = map[State]string{
		Ascii:             "Ascii",
		Unicode:           "Unicode",
		Custom:            "Custom",
		AsciiAz:           "AZ",
		AsciiNums:         "0-9",
		AsciiSpec:         "!$",
		AsciiAll:          "All",
		UnicodeFull:       "█",
		UnicodeHalf:       "▀▄",
		UnicodeQuart:      "▞▟",
		UnicodeShadeLight: "░",
		UnicodeShadeMed:   "▒",
		UnicodeShadeHeavy: "▓",
		OneColor:          "1 Color",
		TwoColor:          "2 Colors",
	}

	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
	titleStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawCharControls() string {
	if m.mode == Custom {
		content := m.drawCustomControls()
		return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Left).Render(content)
	}

	whitespace := 0

	var buttonOrder []State
	switch m.charButtons {
	case Ascii:
		buttonOrder = asciiButtonOrder
	case Unicode:
		buttonOrder = unicodeButtonOrder
	}

	buttons := make([]string, len(buttonOrder))
	for i, state := range buttonOrder {
		buttonStyle := style.NormalButtonNode
		if m.IsActive && state == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if state == m.asciiMode || state == m.unicodeMode {
			buttonStyle = style.ActiveButtonNode
		}

		buttons[i] = buttonStyle.Copy().Render(stateNames[state])

		whitespace += lipgloss.Width(buttons[i])
	}

	gapSpace := whitespace / (len(buttons))
	for i, button := range buttons {
		buttons[i] = lipgloss.NewStyle().PaddingRight(gapSpace).Render(button)
	}
	content := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Left).Render(content)
}

func (m Model) drawCustomControls() string {
	nodeStyle := style.NormalButtonNode.Copy().PaddingRight(1)
	if m.customInput.Focused() {
		nodeStyle = style.ActiveButtonNode.Copy().PaddingRight(1)
	} else if m.focus == SymbolsForm {
		nodeStyle = style.FocusButtonNode.Copy().PaddingRight(1)
	}
	m.customInput.PromptStyle = nodeStyle.Copy()
	return m.customInput.View()
}

func (m Model) drawColorsButtons() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Background Color:")

	yesStyle := style.NormalButtonNode
	if m.IsActive && TwoColor == m.focus {
		yesStyle = style.FocusButtonNode
	} else if m.useFgBg == TwoColor {
		yesStyle = style.ActiveButtonNode
	}
	yesButton := yesStyle.Render("Yes")
	yesButton = lipgloss.NewStyle().Width(5).AlignHorizontal(lipgloss.Center).Render(yesButton)

	noStyle := style.NormalButtonNode
	if m.IsActive && OneColor == m.focus {
		noStyle = style.FocusButtonNode
	} else if m.useFgBg == OneColor {
		noStyle = style.ActiveButtonNode
	}
	noButton := noStyle.Render("No")
	noButton = lipgloss.NewStyle().Width(5).AlignHorizontal(lipgloss.Center).Render(noButton)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, yesButton, noButton)
}
