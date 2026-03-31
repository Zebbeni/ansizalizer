package loader

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

	IsSelected    bool // true if we've selected something (ie. render w/ loader)
	ShouldUnfocus bool

	width int
}

func New(w int) Model {
	fileBrowser := browser.New(paletteExtensions, w-2)

	return Model{
		FileBrowser:   fileBrowser,
		IsSelected:    false,
		ShouldUnfocus: false,
		width:         w,
	}
}

func NewAtDir(dir string, w int) Model {
	m := New(w)
	m.FileBrowser = browser.NewAtDir(dir, paletteExtensions, w-2)
	return m
}

// NewWithPalette returns a Model with a pre-set palette for testing.
func NewWithPalette(colors color.Palette, w int) Model {
	return NewWithNamedPalette("test", colors, w)
}

// NewWithNamedPalette returns a Model with a pre-set named palette.
func NewWithNamedPalette(name string, colors color.Palette, w int) Model {
	m := New(w)
	m.palette = palette.New(name, colors, w-5, 3)
	m.IsSelected = true
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	m.FileBrowser, cmd = m.FileBrowser.Update(msg)

	if m.FileBrowser.SelectedFile != "" && m.FileBrowser.SelectedFile != m.paletteFilepath {
		m.paletteFilepath = m.FileBrowser.SelectedFile

		name := strings.Split(filepath.Base(m.paletteFilepath), ".hex")[0]
		colors, err := parsePaletteFile(m.paletteFilepath)
		if err != nil {
			return m, tea.Batch(cmd, event.BuildDisplayCmd("error parsing paletteFilepath file"))
		}
		m.palette = palette.New(name, colors, m.width-5, 3)

		m.IsSelected = true
		return m, tea.Batch(cmd, event.StartRenderToViewCmd)
	}

	if m.FileBrowser.ShouldClose {
		m.IsSelected = false
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

	title := m.drawTitle()
	browser := m.FileBrowser.View()
	return lipgloss.JoinVertical(lipgloss.Top, title, browser, activePreview)
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
