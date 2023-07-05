package loader

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
	"github.com/Zebbeni/ansizalizer/style"
)

var (
	paletteExtensions = map[string]bool{".hex": true}
)

type Model struct {
	FileBrowser browser.Model

	paletteFilepath string
	palette         palette.Model

	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	fileBrowser := browser.New(paletteExtensions, drawBrowserTitle(), w-2)

	return Model{
		FileBrowser:   fileBrowser,
		ShouldUnfocus: false,
		width:         w,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ActiveFile != m.paletteFilepath {
		m.paletteFilepath = m.FileBrowser.ActiveFile

		name := strings.Split(filepath.Base(m.paletteFilepath), ".hex")[0]
		colors, err := parsePaletteFile(m.paletteFilepath)
		if err != nil {
			return m, tea.Batch(cmd, event.BuildDisplayCmd("error parsing paletteFilepath file"))
		}
		m.palette = palette.New(name, colors, m.width-5, 3)

		return m, tea.Batch(cmd, event.StartRenderToViewCmd)
	}

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
		m.ShouldUnfocus = true
	}

	return m, cmd
}

func (m Model) View() string {
	activePreview := style.DimmedTitle.Render("No palette selected")
	if len(m.palette.Colors()) != 0 {
		activePreview = m.palette.View()
	}
	activePreview = lipgloss.NewStyle().Padding(0, 0, 1, 2).Render(activePreview)

	browser := m.FileBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Top, browser, activePreview)
}

func (m Model) GetCurrent() palette.Model {
	return m.palette
}

func parsePaletteFile(filepath string) (color.Palette, error) {
	readFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var col colorful.Color
	p := make(color.Palette, 0, 256)

	for fileScanner.Scan() {
		col, err = colorful.Hex(fmt.Sprintf("#%s", fileScanner.Text()))
		if err != nil {
			return nil, err
		}
		p = append(p, col)
	}

	return p, nil
}
