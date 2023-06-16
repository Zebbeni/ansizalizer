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
	"github.com/Zebbeni/ansizalizer/controls/settings/colors/description"
	"github.com/Zebbeni/ansizalizer/io"
	"github.com/Zebbeni/ansizalizer/style"
)

var (
	paletteExtensions = []string{".hex"}
)

type Model struct {
	FileBrowser browser.Model

	paletteFilepath string
	palette         color.Palette

	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	fileBrowser := browser.New(paletteExtensions, w-2)

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
	var err error

	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ActiveFile != m.paletteFilepath {
		m.paletteFilepath = m.FileBrowser.ActiveFile
		m.palette, err = parsePaletteFile(m.paletteFilepath)
		if err != nil {
			return m, tea.Batch(cmd, io.BuildDisplayCmd("error parsing paletteFilepath file"))
		}
		return m, tea.Batch(cmd, io.StartRenderCmd)
	}

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
		m.ShouldUnfocus = true
	}

	return m, cmd
}

func (m Model) View() string {
	var activePreview string

	if shouldParse(m.paletteFilepath) {
		name := strings.Split(filepath.Base(m.paletteFilepath), ".hex")[0]
		title := style.SelectedTitle.Render(name)

		palette := m.palette
		desc := description.Palette(palette, m.width-5, 3)

		activePreview = lipgloss.JoinVertical(lipgloss.Top, title, desc)

	} else {
		activePreview = style.DimmedTitle.Render("None")
	}

	activePreview = lipgloss.NewStyle().PaddingBottom(1).PaddingLeft(2).Render(activePreview)

	browser := m.FileBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Top, browser, activePreview)
}

func (m Model) GetCurrent() color.Palette {
	return m.palette
}

func shouldParse(filename string) bool {
	return strings.Contains(filename, ".hex")
}

func parsePaletteFile(filepath string) (color.Palette, error) {
	readFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var col colorful.Color
	palette := make(color.Palette, 0, 256)

	for fileScanner.Scan() {
		col, err = colorful.Hex(fmt.Sprintf("#%s", fileScanner.Text()))
		palette = append(palette, col)
		if err != nil {
			return nil, err
		}
	}

	return palette, nil
}
