package process

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	AlphaPlaceholder string = " "
)

func (m Renderer) outputStrings(rows ...string) (string) {
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}