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
	"github.com/Zebbeni/ansizalizer/style"
)

var (
	paletteExtensions = []string{".hex"}
)

type Model struct {
	FileBrowser browser.Model

	name    string
	palette color.Palette

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
	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.ShouldClose {
		m.FileBrowser.ShouldClose = false
		m.ShouldUnfocus = true
	}
	return m, cmd
}

func (m Model) View() string {
	selected := style.DimmedTitle.Render("None")

	activeFilename := filepath.Base(m.FileBrowser.ActiveFile)
	if shouldParse(activeFilename) {
		name := strings.Split(activeFilename, ".hex")[0]
		selected = style.SelectedTitle.Render(name)

		palette, err := parsePaletteFile(m.FileBrowser.ActiveFile)
		if err == nil {
			selected = lipgloss.JoinVertical(lipgloss.Top, selected, description.Palette(palette, m.width-4, 3))
		}
		selected = lipgloss.NewStyle().PaddingBottom(1).PaddingLeft(2).Render(selected)
	}

	browser := m.FileBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Top, selected, browser)
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
		fmt.Println(err)
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
