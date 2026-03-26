package characters

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder         = []State{Ascii, Unicode, Custom}
	asciiButtonOrder   = []State{AsciiAz, AsciiNums, AsciiSpec, AsciiAll}
	unicodeButtonOrder = []State{UnicodeFull, UnicodeHalf, UnicodeQuart, UnicodeShadeLight, UnicodeShadeMed, UnicodeShadeHeavy}

	selectionModeOrder = []State{DarkToLight, Sequence, Random}

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
		ColorBgOff:        "FG",
		ColorBgOn:         "FG + BG",
		DarkToLight:       "Brightness",
		Sequence:          "Repeat",
		Random:            "Random",
	}
)

func (m Model) drawCharControls() string {
	if m.charControls == Custom {
		return m.drawCustomControls()
	}

	var buttonOrder []State
	switch m.charControls {
	case Ascii:
		buttonOrder = asciiButtonOrder
	case Unicode:
		buttonOrder = unicodeButtonOrder
	}

	innerWidth := m.width - 4
	buttonWidth := innerWidth / len(buttonOrder)

	buttons := make([]string, len(buttonOrder))
	for i, state := range buttonOrder {
		buttonStyle := style.NormalButtonNode
		if m.IsActive && state == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if state == m.asciiMode || state == m.unicodeMode {
			buttonStyle = style.ActiveButtonNode
		}

		buttons[i] = lipgloss.NewStyle().Width(buttonWidth).Render(buttonStyle.Copy().Render(stateNames[state]))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m Model) drawCustomControls() string {
	nodeStyle := style.NormalButtonNode.Copy().PaddingRight(1)
	textStyle := style.BgStyle().Foreground(style.DimmedColor1)
	if m.customInput.Focused() {
		nodeStyle = style.ActiveButtonNode.Copy().PaddingRight(1)
		textStyle = style.BgStyle().Foreground(style.SelectedColor1)
		m.customInput.Cursor.SetMode(cursor.CursorBlink)
	} else if m.focus == SymbolsForm {
		nodeStyle = style.FocusButtonNode.Copy().PaddingRight(1)
		textStyle = style.BgStyle().Foreground(style.NormalColor1)
		m.customInput.Cursor.SetMode(cursor.CursorHide)
	} else {
		m.customInput.Cursor.SetMode(cursor.CursorHide)
	}
	m.customInput.PromptStyle = nodeStyle
	m.customInput.TextStyle = textStyle

	symbolsRow := m.customInput.View()
	selectionRow := style.BgStyle().PaddingTop(1).Render(m.drawSelectionMode())
	return lipgloss.JoinVertical(lipgloss.Center, symbolsRow, selectionRow)
}

var modeDescriptions = map[State]string{
	DarkToLight: "Map chars by brightness",
	Sequence:    "Repeat chars in order",
	Random:      "Assign chars randomly",
}

func (m Model) drawSelectionMode() string {
	buttons := make([]string, len(selectionModeOrder))
	for i, state := range selectionModeOrder {
		buttonStyle := style.NormalButtonNode
		if m.IsActive && state == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if state == m.selectionMode {
			buttonStyle = style.ActiveButtonNode
		}
		buttons[i] = buttonStyle.Render(stateNames[state])
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	// Show description for focused mode
	if desc, ok := modeDescriptions[m.focus]; ok {
		descInner := style.DimmedTitle.Copy().Italic(true).Render(desc)
		descText := lipgloss.NewStyle().PaddingTop(1).Width(m.width - 4).AlignHorizontal(lipgloss.Center).Render(descInner)
		return lipgloss.JoinVertical(lipgloss.Left, row, descText)
	}

	return row
}

func (m Model) drawColorsButtons() string {
	title := style.DimmedTitle.Copy().PaddingLeft(1).Render("Colors:")

	oneStyle := style.NormalButtonNode
	if m.IsActive && ColorBgOff == m.focus {
		oneStyle = style.FocusButtonNode
	} else if !m.colorBg {
		oneStyle = style.ActiveButtonNode
	}
	oneButton := oneStyle.Render("FG")
	oneButton = style.BgStyle().Width(5).AlignHorizontal(lipgloss.Center).Render(oneButton)

	twoStyle := style.NormalButtonNode
	if m.IsActive && ColorBgOn == m.focus {
		twoStyle = style.FocusButtonNode
	} else if m.colorBg {
		twoStyle = style.ActiveButtonNode
	}
	twoButton := twoStyle.Render("FG + BG")
	twoButton = style.BgStyle().Width(9).AlignHorizontal(lipgloss.Center).Render(twoButton)

	return lipgloss.JoinHorizontal(lipgloss.Left, title, oneButton, twoButton)
}
