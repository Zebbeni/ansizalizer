package characters

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder         = []State{Ascii, Unicode, Custom}
	asciiButtonOrder   = []State{AsciiAz, AsciiNums, AsciiSpec, AsciiAll}
	unicodeButtonOrder = []State{UnicodeFull, UnicodeHalf, UnicodeQuart, UnicodeShadeLight, UnicodeShadeMed, UnicodeShadeHeavy, UnicodeShadeAll}

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
		UnicodeShadeAll:   "░▒▓",
		OneColor:          "1 Color",
		TwoColor:          "2 Colors",
	}

	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
	titleStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawCharButtons() string {
	var buttonOrder []State
	switch m.charButtons {
	case Ascii:
		buttonOrder = asciiButtonOrder
	case Unicode:
		buttonOrder = unicodeButtonOrder
	}
	buttons := make([]string, len(buttonOrder))
	for i, state := range buttonOrder {
		styleColor := normalColor
		if m.IsActive && state == m.focus {
			styleColor = focusColor
		} else if state == m.active {
			styleColor = activeColor
		}
		style := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styleColor).
			Foreground(styleColor)

		// quick dirty stuff to make buttons fit nicely. Let's do this in a smarter / cleaner way later
		if m.charButtons == Unicode {
			buttons[i] = style.Copy().Render(stateNames[state])
		} else {
			buttons[i] = style.Copy().Padding(0, 0).Render(stateNames[state])
		}
	}
	content := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	return lipgloss.NewStyle().Width(m.width).AlignHorizontal(lipgloss.Center).Render(content)
}

func (m Model) drawColorsButtons() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Use Background:")

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
