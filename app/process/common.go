package process

import (
	"bufio"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	AlphaPlaceholder string = " "
)

// ansiColor formats ANSI true-color escape sequences with correct RGB values,
// bypassing a bug in termenv where uint8(f.R*255) truncates instead of
// rounding, making some RGB values (like 98) unreachable via hex input.
func ansiFg(c colorful.Color) string {
	r := uint8(c.R*255.0 + 0.5)
	g := uint8(c.G*255.0 + 0.5)
	b := uint8(c.B*255.0 + 0.5)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func ansiBg(c colorful.Color) string {
	r := uint8(c.R*255.0 + 0.5)
	g := uint8(c.G*255.0 + 0.5)
	b := uint8(c.B*255.0 + 0.5)
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

const ansiReset = "\x1b[0m"

func renderChar(char string, fg colorful.Color, useBg bool, bg colorful.Color, bold bool) string {
	s := ansiFg(fg)
	if useBg {
		s += ansiBg(bg)
	}
	if bold {
		s += "\x1b[1m"
	}
	return s + char + ansiReset
}


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