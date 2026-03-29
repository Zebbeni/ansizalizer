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
	prompt := style.DimmedTitle.Render("Dithering:")
	prompt = style.BgStyle().Width(12).Render(prompt)

	nodeStyle := style.NormalButtonNode
	if m.IsActive && m.focus == DitherOn {
		nodeStyle = style.FocusButtonNode
	} else if m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	onNode := style.BgStyle().Width(4).Render(nodeStyle.Copy().Render("On"))

	nodeStyle = style.NormalButtonNode
	if m.IsActive && m.focus == DitherOff {
		nodeStyle = style.FocusButtonNode
	} else if !m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	offNode := nodeStyle.Copy().Render("Off")

	return lipgloss.JoinHorizontal(lipgloss.Left, prompt, onNode, offNode)
}

func (m Model) drawSerpentineOptions() string {
	prompt := style.DimmedTitle.Render("Serpentine:")
	prompt = style.BgStyle().Width(12).Render(prompt)

	nodeStyle := style.NormalButtonNode
	if m.IsActive && m.focus == SerpentineOn {
		nodeStyle = style.FocusButtonNode
	} else if m.doSerpentine {
		nodeStyle = style.ActiveButtonNode
	}
	onNode := style.BgStyle().Width(4).Render(nodeStyle.Copy().Render("On"))

	nodeStyle = style.NormalButtonNode
	if m.IsActive && m.focus == SerpentineOff {
		nodeStyle = style.FocusButtonNode
	} else if !m.doSerpentine {
		nodeStyle = style.ActiveButtonNode
	}
	offNode := nodeStyle.Copy().Render("Off")

	return lipgloss.JoinHorizontal(lipgloss.Left, prompt, onNode, offNode)
}

func (m Model) drawModeTabs() string {
	tabs := []style.Tab{
		{Label: modeTabNames[ModeMatrix], Focused: m.focus == ModeMatrix, Active: m.modeControls == ModeMatrix},
		{Label: modeTabNames[ModeBayer], Focused: m.focus == ModeBayer, Active: m.modeControls == ModeBayer},
		{Label: modeTabNames[ModeClusteredDot], Focused: m.focus == ModeClusteredDot, Active: m.modeControls == ModeClusteredDot},
	}
	return style.RenderTabs(tabs, m.IsActive, m.drawModeContent(), m.width)
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
	promptStyle := style.DimmedTitle.Copy()
	textStyle := style.BgStyle().Foreground(style.DimmedColor1)
	if m.strengthInput.Focused() {
		promptStyle = style.SelectedTitle.Copy()
		textStyle = style.BgStyle().Foreground(style.SelectedColor1)
	} else if m.focus == StrengthForm && m.IsActive {
		promptStyle = style.NormalTitle.Copy()
		textStyle = style.BgStyle().Foreground(style.NormalColor1)
	}

	styles := m.strengthInput.Styles()
	styles.Focused.Prompt = promptStyle
	styles.Focused.Text = textStyle
	styles.Blurred.Prompt = promptStyle
	styles.Blurred.Text = textStyle
	styles.Cursor.Blink = m.strengthInput.Focused()
	m.strengthInput.SetStyles(styles)
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
