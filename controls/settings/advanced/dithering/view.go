package dithering

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle  = style.BgStyle().Border(inactiveTabBorder, true)
	activeTabStyle    = style.BgStyle().Border(activeTabBorder, true)
	windowStyle       = style.BgStyle().Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop().Padding(1, 0)

	modeTabNames = map[State]string{
		ModeMatrix:       "Preset",
		ModeBayer:        "Bayer",
		ModeClusteredDot: "Dot",
	}
)

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
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []State{ModeMatrix, ModeBayer, ModeClusteredDot}

	borderColor := style.DimmedColor2
	if m.IsActive {
		borderColor = style.NormalColor1
	}

	focusIsTab := m.focus == ModeMatrix || m.focus == ModeBayer || m.focus == ModeClusteredDot
	focusedTabIsActive := focusIsTab && m.focus == m.modeControls
	focusIsContent := m.IsActive && !focusIsTab
	useDoubleContent := m.IsActive && (focusedTabIsActive || focusIsContent)

	for i, t := range tabs {
		isFirst := i == 0
		isLast := i == len(tabs)-1
		isFocused := m.focus == t
		showControls := m.modeControls == t

		fgColor := style.DimmedColor2
		if m.IsActive {
			if isFocused {
				fgColor = style.SelectedColor1
			} else {
				fgColor = style.DimmedColor1
			}
		} else {
			if isFocused {
				fgColor = style.NormalColor2
			}
		}

		var tabStyle lipgloss.Style
		if m.IsActive && isFocused && showControls {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if m.IsActive && isFocused {
			tabStyle = style.FocusTabStyle.Copy()
		} else if showControls && useDoubleContent {
			tabStyle = style.FocusActiveTabStyle.Copy()
		} else if showControls {
			tabStyle = style.ActiveTabStyle.Copy()
		} else if useDoubleContent {
			tabStyle = style.InactiveTabOnHeavyStyle.Copy()
		} else {
			tabStyle = style.InactiveTabStyle.Copy()
		}

		border, _, _, _, _ := tabStyle.GetBorder()
		if useDoubleContent {
			if isFirst && showControls {
				border.BottomLeft = "┃"
			} else if isFirst {
				border.BottomLeft = "┢"
			}
			if isLast && showControls {
				border.BottomRight = "┗"
			} else if isLast && isFocused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┷"
			}
		} else {
			if isFirst && showControls {
				border.BottomLeft = "│"
			} else if isFirst && m.IsActive && isFocused {
				border.BottomLeft = "┞"
			} else if isFirst {
				border.BottomLeft = "├"
			}
			if isLast && showControls {
				border.BottomRight = "└"
			} else if isLast && m.IsActive && isFocused {
				border.BottomRight = "┸"
			} else if isLast {
				border.BottomRight = "┴"
			}
		}

		tabStyle = tabStyle.Border(border).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Foreground(fgColor)
		renderedTabs = append(renderedTabs, tabStyle.Render(modeTabNames[t]))
	}

	tabBlock := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	extW := max(m.width-lipgloss.Width(tabBlock)-4, 0)

	var extBorder lipgloss.Border
	if useDoubleContent {
		extBorder = lipgloss.Border{BottomLeft: "━", Bottom: "━", BottomRight: "┓"}
	} else {
		extBorder = lipgloss.Border{BottomLeft: "─", Bottom: "─", BottomRight: "┐"}
	}
	extendedStyle := style.TabWindowStyle.Copy().Border(extBorder).BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Padding(0)
	extended := extendedStyle.Copy().Width(extW).Height(1).Render("")
	renderedTabs = append(renderedTabs, extended)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	winStyle := style.TabWindowStyle
	if useDoubleContent {
		winStyle = style.TabWindowHeavyStyle
	}
	content := m.drawModeContent()
	doc.WriteString(winStyle.Copy().BorderForeground(borderColor).BorderBackground(style.ActiveTheme.Bg).Width(lipgloss.Width(row)).Render(content))
	return doc.String()
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
