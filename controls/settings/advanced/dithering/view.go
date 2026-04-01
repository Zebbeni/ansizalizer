package dithering

import (
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var modeTabNames = map[State]string{
	ModeMatrix:       "Preset",
	ModeBayer:        "Bayer",
	ModeClusteredDot: "Dot",
}

func (m Model) drawDitheringOptions() string {
	focused := m.IsActive && (m.focus == DitherOn || m.focus == DitherOff)
	return style.RenderCheckboxFixedWidth("Dithering", m.doDithering, focused, 11)
}

func (m Model) drawSerpentineOptions() string {
	focused := m.IsActive && (m.focus == SerpentineOn || m.focus == SerpentineOff)
	return style.RenderCheckboxFixedWidth("Serpentine", m.doSerpentine, focused, 11)
}

var tabStates = map[State]bool{
	ModeMatrix: true, ModeBayer: true, ModeClusteredDot: true,
	MatrixList: true, ClusteredDotList: true,
	BayerSize2: true, BayerSize4: true, BayerSize8: true, BayerSize16: true,
}

func (m Model) drawModeTabs() string {
	tabs := []style.Tab{
		{Label: modeTabNames[ModeMatrix], Focused: m.focus == ModeMatrix, Active: m.modeControls == ModeMatrix},
		{Label: modeTabNames[ModeBayer], Focused: m.focus == ModeBayer, Active: m.modeControls == ModeBayer},
		{Label: modeTabNames[ModeClusteredDot], Focused: m.focus == ModeClusteredDot, Active: m.modeControls == ModeClusteredDot},
	}
	// Only pass isActive when focus is on a tab or tab content
	tabActive := m.IsActive && tabStates[m.focus]
	return style.RenderTabs(tabs, tabActive, m.drawModeContent(), m.width)
}

func (m Model) drawModeContent() string {
	switch m.modeControls {
	case ModeMatrix:
		return m.matrixList.View()
	case ModeBayer:
		return m.drawBayerSizeButtons()
	case ModeClusteredDot:
		return m.clusteredDotList.View()
	}
	return ""
}

func (m Model) drawBayerSizeButtons() string {
	prompt := style.DimmedTitle.Render("Size:")

	type sizeOption struct {
		state State
		size  uint
		label string
	}
	options := []sizeOption{
		{BayerSize2, 2, "2"},
		{BayerSize4, 4, "4"},
		{BayerSize8, 8, "8"},
		{BayerSize16, 16, "16"},
	}

	buttons := make([]string, len(options))
	for i, opt := range options {
		btnStyle := style.NormalButtonNode
		if m.IsActive && m.focus == opt.state {
			btnStyle = style.FocusButtonNode
		} else if m.bayerSize == opt.size {
			btnStyle = style.ActiveButtonNode
		}
		buttons[i] = style.BgStyle().Width(4).Render(btnStyle.Render(opt.label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, append([]string{prompt}, buttons...)...)
}

func (m Model) drawStrength() string {
	promptStyle := style.DimmedTitle
	textStyle := style.BgStyle().Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if m.strengthInput.Focused() {
		promptStyle = style.SelectedTitle
		textStyle = style.BgStyle().Foreground(style.SelectedColor1)
		cursorColor = style.SelectedColor1
	} else if m.focus == StrengthForm && m.IsActive {
		promptStyle = style.NormalTitle
		textStyle = style.BgStyle().Foreground(style.NormalColor1)
		cursorColor = style.NormalColor1
	}

	styles := m.strengthInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textStyle
	styles.Blurred.Prompt = promptStyle
	styles.Blurred.Text = textStyle
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.strengthInput.SetStyles(styles)
	m.strengthInput.SetVirtualCursor(m.strengthInput.Focused())
	return m.strengthInput.View()
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
