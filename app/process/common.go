package process

import (
	"strings"
	"bufio"
	"github.com/charmbracelet/lipgloss"
)

const (
	AlphaPlaceholder string = " "
)

func (m Renderer) outputStrings(rows ...string) (string) {
	content := ""
	if m.Settings.Alpha.ShouldOutputAlpha() && m.Settings.Alpha.TrimAlpha() {
		contentAlpha := lipgloss.JoinVertical(lipgloss.Left, rows...)
		// iterate through the return of JoinVertical, separating by lines, trimming our AlphaPlaceholder from the right, and then recombining
		reader := strings.NewReader(contentAlpha)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			content += strings.TrimRight(scanner.Text(), " ") + "\n"
		}
		content = strings.TrimRight(content, "\n")
	} else {
		content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	}
	return content
}