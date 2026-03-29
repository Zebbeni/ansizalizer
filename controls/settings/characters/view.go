package characters

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder         = []State{Ascii, Unicode, Custom}
	asciiButtonOrder   = []State{AsciiAz, AsciiNums, AsciiSpec, AsciiAll}
	unicodeButtonOrder = []State{UnicodeFull, UnicodeHalf, UnicodeQuart, UnicodeShadeLight, UnicodeShadeMed, UnicodeShadeHeavy}

	selectionModeOrder = []State{Variance, Sequence, Random}

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
		Variance:          "Variance",
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
	row := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	if m.charControls == Ascii && m.colorBg {
		thresholdRow := style.BgStyle().PaddingTop(1).Render(m.drawThresholdInput())
		return lipgloss.JoinVertical(lipgloss.Left, row, thresholdRow)
	}

	return row
}

func (m Model) drawCustomControls() string {
	nodeStyle := style.NormalButtonNode.Copy().PaddingRight(1)
	textStyle := style.BgStyle().Foreground(style.DimmedColor1)
	cursorBlink := false
	if m.customInput.Focused() {
		nodeStyle = style.ActiveButtonNode.Copy().PaddingRight(1)
		textStyle = style.BgStyle().Foreground(style.SelectedColor1)
		cursorBlink = true
	} else if m.focus == SymbolsForm {
		nodeStyle = style.FocusButtonNode.Copy().PaddingRight(1)
		textStyle = style.BgStyle().Foreground(style.NormalColor1)
	}
	styles := m.customInput.Styles()
	styles.Focused.Prompt = nodeStyle
	styles.Focused.Text = textStyle
	styles.Cursor.Blink = cursorBlink
	m.customInput.SetStyles(styles)

	symbolsRow := m.customInput.View()
	selectionRow := style.BgStyle().PaddingTop(1).Render(m.drawSelectionMode())

	if m.colorBg && m.selectionMode == Variance {
		thresholdRow := style.BgStyle().PaddingTop(1).Render(m.drawThresholdInput())
		return lipgloss.JoinVertical(lipgloss.Center, symbolsRow, selectionRow, thresholdRow)
	}

	return lipgloss.JoinVertical(lipgloss.Center, symbolsRow, selectionRow)
}

var modeDescriptions = map[State]string{
	Variance: "Select by color variance left: less, right: more",
	Sequence: "Repeat chars in order",
	Random:   "Assign chars randomly",
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

	// Show description for focused mode or when threshold is focused
	showFocus := m.focus
	if m.focus == ThresholdForm || m.thresholdInput.Focused() {
		showFocus = Variance
	}
	if desc, ok := modeDescriptions[showFocus]; ok {
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

func (m Model) drawThresholdInput() string {
	promptStyle := style.DimmedTitle.Copy()
	textStyle := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	cursorBlink := false
	if m.thresholdInput.Focused() {
		promptStyle = style.SelectedTitle.Copy()
		textStyle = lipgloss.NewStyle().Foreground(style.SelectedColor1)
		cursorBlink = true
	} else if m.focus == ThresholdForm && m.IsActive {
		promptStyle = style.NormalTitle.Copy()
		textStyle = lipgloss.NewStyle().Foreground(style.NormalColor1)
	}
	styles := m.thresholdInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textStyle
	styles.Cursor.Blink = cursorBlink
	m.thresholdInput.SetStyles(styles)
	thresholdRow := m.thresholdInput.View()

	if m.focus == ThresholdForm || m.thresholdInput.Focused() {
		desc := style.DimmedTitle.Copy().Italic(true).Render("Variance threshold for displaying characters.")
		descRow := lipgloss.NewStyle().PaddingLeft(1).PaddingTop(1).Width(m.width - 4).Render(desc)
		return lipgloss.JoinVertical(lipgloss.Left, thresholdRow, descRow)
	}

	return thresholdRow
}
