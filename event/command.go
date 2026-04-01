package event

import (
	"fmt"
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
)

// CheckThemeMsg is sent when a component changes something that may affect the app theme.
// The top-level app compares against its current state and applies changes if needed.
type CheckThemeMsg struct {
	Palette  color.Palette
	ThemeName string
}

// BuildCheckThemeCmd returns a command that sends a CheckThemeMsg with the given palette and theme name.
func BuildCheckThemeCmd(palette color.Palette, themeName string) tea.Cmd {
	return func() tea.Msg {
		return CheckThemeMsg{Palette: palette, ThemeName: themeName}
	}
}

// StartRenderToViewMsg signals the app to begin rendering the current image for the viewer.
type StartRenderToViewMsg bool

// StartRenderToViewCmd is a tea.Cmd that sends a StartRenderToViewMsg to initiate rendering.
func StartRenderToViewCmd() tea.Msg {
	return StartRenderToViewMsg(true)
}

// FinishRenderToViewMsg carries the result of a single-image render to the viewer.
type FinishRenderToViewMsg struct {
	FilePath     string
	ImgString    string
	ColorsString string
	AlphaString  string
	ExtraInfo    string
}

// FinishRenderGIFToViewMsg carries the result of a GIF render, including all frames and delay.
type FinishRenderGIFToViewMsg struct {
	FilePath     string
	Frames       []string
	Delay        time.Duration
	ColorsString string
	AlphaString  string
	ExtraInfo    string
}

// AnimationTickMsg triggers the viewer to advance to the next GIF frame.
type AnimationTickMsg struct {
	Generation int
}

// BuildAnimationTickCmd returns a command that sends an AnimationTickMsg after the given delay.
// The generation parameter ties the tick to a specific animation sequence, preventing stale ticks.
func BuildAnimationTickCmd(delay time.Duration, generation int) tea.Cmd {
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return AnimationTickMsg{Generation: generation}
	})
}

// StartRenderToExportMsg signals the app to begin rendering the current image for export.
type StartRenderToExportMsg bool

// StartRenderToExportCmd is a tea.Cmd that sends a StartRenderToExportMsg to initiate export rendering.
func StartRenderToExportCmd() tea.Msg {
	return StartRenderToExportMsg(true)
}

// FinishRenderToExportMsg carries the result of a render destined for file export.
type FinishRenderToExportMsg struct {
	FilePath     string
	ImgString    string
	ColorsString string
	AlphaString  string
}

// BuildFinishRenderToExportCmd returns a command that delivers a completed export render result.
func BuildFinishRenderToExportCmd(msg FinishRenderToExportMsg) tea.Cmd {
	return func() tea.Msg { return msg }
}

// StartAdaptingMsg signals the app to begin generating an adaptive color palette from the current image.
type StartAdaptingMsg bool

// StartAdaptingCmd is a tea.Cmd that sends a StartAdaptingMsg to initiate palette generation.
func StartAdaptingCmd() tea.Msg {
	return StartAdaptingMsg(true)
}

// FinishAdaptingMsg carries the result of an adaptive palette generation, including the palette name and colors.
type FinishAdaptingMsg struct {
	Name   string
	Colors color.Palette
}

// StartExportMsg initiates a batch export operation with the specified source and destination paths.
type StartExportMsg struct {
	SourcePath      string
	DestinationPath string
	IsDir           bool
	UseSubDirs      bool
}

// BuildStartExportCmd returns a command that sends the given StartExportMsg.
func BuildStartExportCmd(msg StartExportMsg) tea.Cmd {
	return func() tea.Msg { return msg }
}

// FinishExportMsg signals that a batch export operation has completed.
type FinishExportMsg bool

// FinishExportingCmd is a tea.Cmd that sends a FinishExportMsg.
func FinishExportingCmd() tea.Msg {
	return FinishExportMsg(true)
}

// DisplayMsg could eventually contain a type
// that indicates what style to use (warning, error, etc.)
type DisplayMsg string

// BuildDisplayCmd returns a command that sends a DisplayMsg with the given text to the status bar.
func BuildDisplayCmd(msg string) tea.Cmd {
	return func() tea.Msg { return DisplayMsg(msg) }
}

// ClearDisplayCmd is a tea.Cmd that sends an empty DisplayMsg to clear the status bar.
func ClearDisplayCmd() tea.Msg {
	return DisplayMsg("")
}

// LospecRequestMsg represents a paginated request to the Lospec.com palette API.
type LospecRequestMsg struct {
	ID   int
	Page int
	URL  string
}

// BuildLospecRequestCmd returns a command that sends a LospecRequestMsg and displays a loading status.
func BuildLospecRequestCmd(msg LospecRequestMsg) tea.Cmd {
	display := fmt.Sprintf("loading palettes")
	return tea.Batch(func() tea.Msg { return msg }, BuildDisplayCmd(display))
}

// LospecData holds the JSON response from the Lospec.com palette list API.
type LospecData struct {
	Palettes []struct {
		Colors []string `json:"colors"`
		Title  string   `json:"title"`
	} `json:"palettes"`
	TotalCount int `json:"totalCount"`
}

// LospecResponseMsg carries a parsed Lospec API response back to the palette browser.
type LospecResponseMsg struct {
	ID   int
	Page int
	Data LospecData
}

// BuildLospecResponseCmd returns a command that delivers a LospecResponseMsg and clears the status bar.
func BuildLospecResponseCmd(msg LospecResponseMsg) tea.Cmd {
	return tea.Batch(func() tea.Msg { return msg }, ClearDisplayCmd)
}
