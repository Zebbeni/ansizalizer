package process

import (
	"strings"
	"bufio"
	"regexp"
	"math"
	"github.com/charmbracelet/lipgloss"
)

const (
	AlphaPlaceholder string = " "
)

func (m Renderer) outputStrings(rows ...string) (string) {
	content := ""
	if m.Settings.Alpha.ShouldOutputAlpha() && m.Settings.Alpha.TrimAlpha() {
		leftWhitespaceRE := `(^\s+)`
		re := regexp.MustCompile(leftWhitespaceRE)
		leftTrimAmount := math.MaxInt
		contentAlpha := lipgloss.JoinVertical(lipgloss.Left, rows...)
		// get the minimum amount of left whitespace characters we can trim
		reader := strings.NewReader(contentAlpha)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			leftWhitespaceMatch := re.FindStringSubmatch(scanner.Text())
			if leftWhitespaceMatch != nil && leftTrimAmount > len(leftWhitespaceMatch[0]) { // set our minimum trim
				leftTrimAmount = len(leftWhitespaceMatch[0])
			}
			if leftTrimAmount == 0 { // this image runs all the way to the left, so no need to trim
				break
			}
		}
		if leftTrimAmount == math.MaxInt { // hopefully never possible! (this would be a fully blank image)
			leftTrimAmount = 0
		}
		// iterate through the return of JoinVertical, separating by lines, trimming whitespace from the left if necessary, then trimming our AlphaPlaceholder from the right, and then recombining
		blankLine := ([]string)(nil)
		blankLineRE := `(^\s+$)`
		re = regexp.MustCompile(blankLineRE)
		imageTop := true
		reader = strings.NewReader(contentAlpha)
		scanner = bufio.NewScanner(reader)
		for scanner.Scan() {
			thisLine := scanner.Text()
			if (leftTrimAmount > 0) { // trim from the left
				thisLine = thisLine[leftTrimAmount:]
			}
			if imageTop { // we only care about blankLines at the top of the image, not in the middle
				blankLine = re.FindStringSubmatch(thisLine)
			}
			if blankLine == nil { // we have characters, so we're not at the top of the image
				imageTop = false
			}
			if !imageTop {
				content += strings.TrimRight(thisLine, " ") + "\n" // trim from the right
			}
		}
		content = strings.TrimRight(content, "\n") // trim the final carriage return
	} else {
		content += lipgloss.JoinVertical(lipgloss.Left, rows...)
	}
	return content
}