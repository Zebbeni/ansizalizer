package dithering

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) drawDitheringOptions() string {
	prompt := style.DimmedTitle.Render("Use Dithering:")
	prompt = lipgloss.NewStyle().Width(15).Render(prompt)

	nodeStyle := style.NormalButtonNode
	if m.IsActive && m.focus == DitherOn {
		nodeStyle = style.FocusButtonNode
	} else if m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	onNode := lipgloss.NewStyle().Width(4).Render(nodeStyle.Copy().Render("On"))

	nodeStyle = style.NormalButtonNode
	if m.IsActive && m.focus == DitherOff {
		nodeStyle = style.FocusButtonNode
	} else if m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	offNode := nodeStyle.Copy().Render("Off")

	return lipgloss.JoinHorizontal(lipgloss.Left, prompt, onNode, offNode)
}

func (m Model) drawSerpentineOptions() string {
	prompt := style.DimmedTitle.Render("Do Serpentine:")
	prompt = lipgloss.NewStyle().Width(15).Render(prompt)

	nodeStyle := style.NormalButtonNode
	if m.IsActive && m.focus == SerpentineOn {
		nodeStyle = style.FocusButtonNode
	} else if m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	onNode := lipgloss.NewStyle().Width(4).Render(nodeStyle.Copy().Render("On"))

	nodeStyle = style.NormalButtonNode
	if m.IsActive && m.focus == SerpentineOff {
		nodeStyle = style.FocusButtonNode
	} else if m.doDithering {
		nodeStyle = style.ActiveButtonNode
	}
	offNode := nodeStyle.Copy().Render("Off")

	return lipgloss.JoinHorizontal(lipgloss.Left, prompt, onNode, offNode)
}

func (m Model) drawMatrix() string {
	prompt := style.DimmedTitle.Render("Select Matrix")
	return lipgloss.JoinVertical(lipgloss.Left, prompt, m.list.View())
}
