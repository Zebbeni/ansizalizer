package characters

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateOrder         = []State{Unicode, Ascii, Custom}
	asciiButtonOrder   = []State{AsciiAz, AsciiNums, AsciiSpec, AsciiAll}
	unicodeButtonOrder = []State{UnicodeFull, UnicodeHalf, UnicodeQuart, UnicodeShadeLight, UnicodeShadeMed, UnicodeShadeHeavy}

	varianceModes  = []State{DarkVariance, LightVariance}
	sequenceModes  = []State{Sequence, Random}

	stateNames = map[State]string{
		Ascii:             "Ascii",
		Unicode:           "Block",
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
		DarkVariance:      "Dark Var",
		LightVariance:     "Light Var",
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

	buttonWidth := m.width / len(buttonOrder)

	buttons := make([]string, len(buttonOrder))
	for i, state := range buttonOrder {
		buttonStyle := style.NormalButtonNode
		if m.IsActive && state == m.focus {
			buttonStyle = style.FocusButtonNode
		} else if state == m.asciiMode || state == m.unicodeMode {
			buttonStyle = style.ActiveButtonNode
		}

		buttons[i] = lipgloss.NewStyle().Width(buttonWidth).Render(buttonStyle.Render(stateNames[state]))
	}
	row := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)

	centeredRow := lipgloss.NewStyle().Width(m.width).Padding(0, 1).AlignHorizontal(lipgloss.Center).Render(row)

	if m.charControls == Ascii {
		paddedRow := lipgloss.NewStyle().PaddingTop(1).Render(centeredRow)
		selectionRow := style.BgStyle().PaddingTop(1).Width(m.width).AlignHorizontal(lipgloss.Center).Render(m.drawSelectionMode())
		return lipgloss.JoinVertical(lipgloss.Center, paddedRow, selectionRow)
	}

	// Unicode/Block: top and bottom padding
	return lipgloss.NewStyle().Padding(1, 0).Render(centeredRow)
}

func (m Model) inputStyles(state State, focused bool) (lipgloss.Style, lipgloss.Style, color.Color) {
	promptStyle := style.DimmedTitle
	textStyle := lipgloss.NewStyle().Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if focused {
		promptStyle = style.SelectedTitle
		textStyle = lipgloss.NewStyle().Foreground(style.SelectedColor1)
		cursorColor = style.SelectedColor1
	} else if m.focus == state && m.IsActive {
		promptStyle = style.NormalTitle
		textStyle = lipgloss.NewStyle().Foreground(style.NormalColor1)
		cursorColor = style.NormalColor1
	}
	return promptStyle, textStyle, cursorColor
}

func (m Model) drawCustomControls() string {
	promptStyle, textStyle, cursorColor := m.inputStyles(SymbolsForm, m.customInput.Focused())
	styles := m.customInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textStyle
	styles.Blurred.Prompt = promptStyle
	styles.Blurred.Text = textStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.customInput.SetStyles(styles)
	m.customInput.SetVirtualCursor(m.customInput.Focused())

	symbolsRow := lipgloss.NewStyle().PaddingTop(1).Render(m.customInput.View())
	selectionRow := style.BgStyle().PaddingTop(1).Width(m.width - 4).AlignHorizontal(lipgloss.Center).Render(m.drawSelectionMode())

	return lipgloss.JoinVertical(lipgloss.Center, symbolsRow, selectionRow)
}

func (m Model) drawThresholdInput() string {
	promptStyle, textSt, cursorColor := m.inputStyles(ThresholdForm, m.thresholdInput.Focused())
	styles := m.thresholdInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textSt
	styles.Blurred.Prompt = promptStyle
	styles.Blurred.Text = textSt
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.thresholdInput.SetStyles(styles)
	m.thresholdInput.SetVirtualCursor(m.thresholdInput.Focused())
	return m.thresholdInput.View()
}

func (m Model) drawSeedInput() string {
	promptStyle, textSt, cursorColor := m.inputStyles(SeedForm, m.seedInput.Focused())
	styles := m.seedInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textSt
	styles.Blurred.Prompt = promptStyle
	styles.Blurred.Text = textSt
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.seedInput.SetStyles(styles)
	m.seedInput.SetVirtualCursor(m.seedInput.Focused())
	return m.seedInput.View()
}

func (m Model) drawSelectionMode() string {
	renderRow := func(modes []State) string {
		buttons := make([]string, len(modes))
		for i, state := range modes {
			buttonStyle := style.NormalButtonNode
			if m.IsActive && state == m.focus {
				buttonStyle = style.FocusButtonNode
			} else if state == m.selectionMode {
				buttonStyle = style.ActiveButtonNode
			}
			buttons[i] = buttonStyle.Render(stateNames[state])
		}
		return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	}

	row1 := renderRow(varianceModes)
	parts := []string{row1}

	// Show threshold input when a variance mode is selected
	if m.selectionMode == DarkVariance || m.selectionMode == LightVariance {
		parts = append(parts, m.drawThresholdInput())
	}

	row2 := renderRow(sequenceModes)
	parts = append(parts, row2)

	// Show seed input when Random is selected
	if m.selectionMode == Random {
		parts = append(parts, m.drawSeedInput())
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

