package event

import (
	"fmt"
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
)

type StartRenderMsg bool

func StartRenderCmd() tea.Msg {
	return StartRenderMsg(true)
}

type FinishRenderMsg struct {
	FilePath  string
	ImgString string
}

type StartAdaptingMsg bool

func StartAdaptingCmd() tea.Msg {
	return StartAdaptingMsg(true)
}

type FinishAdaptingMsg struct {
	Name   string
	Colors color.Palette
}

// DisplayMsg could eventually contain a type
// that indicates what style to use (warning, error, etc.)
type DisplayMsg string

func BuildDisplayCmd(msg string) tea.Cmd {
	return func() tea.Msg { return DisplayMsg(msg) }
}

func ClearDisplayCmd() tea.Msg {
	return DisplayMsg("")
}

func BuildAdaptingCmd() tea.Cmd {
	return tea.Batch(StartAdaptingCmd, BuildDisplayCmd("Generating palette..."))
}

// LospecRequestMsg is a url request used to get a list of
type LospecRequestMsg struct {
	ID   int
	Page int
	URL  string
}

func BuildLospecRequestCmd(msg LospecRequestMsg) tea.Cmd {
	display := fmt.Sprintf("loading palettes")
	return tea.Batch(func() tea.Msg { return msg }, BuildDisplayCmd(display))
}

type LospecData struct {
	Palettes []struct {
		Colors []string `json:"colors"`
		Title  string   `json:"title"`
	} `json:"palettes"`
	TotalCount int `json:"totalCount"`
}

type LospecResponseMsg struct {
	ID   int
	Page int
	Data LospecData
}

func BuildLospecResponseCmd(msg LospecResponseMsg) tea.Cmd {
	return tea.Batch(func() tea.Msg { return msg }, ClearDisplayCmd)
}
