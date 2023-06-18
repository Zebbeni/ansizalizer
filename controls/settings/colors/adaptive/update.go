package adaptive

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var navMap = map[Direction]map[State]State{
	Right: {CountForm: IterForm},
	Left:  {IterForm: CountForm},
	Up:    {Generate: CountForm, Save: Generate},
	Down:  {CountForm: Generate, IterForm: Generate, Generate: Save},
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	switch m.active {
	case CountForm:
		m.countInput.Focus()
		return m, nil
	case IterForm:
		m.iterInput.Focus()
		return m, nil
	case Save:
		return m.savePaletteFile()
	}
	return m, event.BuildAdaptingCmd()
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			m.focus = next
		} else {
			m.ShouldUnfocus = true
		}
	}

	return m, cmd
}

func (m Model) savePaletteFile() (Model, tea.Cmd) {
	filename := fmt.Sprintf("%s.hex", m.palette.Name())

	f, err := os.Create(filename)

	if err != nil {
		return m, event.BuildDisplayCmd("error saving palette file")
	}

	defer f.Close()

	var hexStrings string

	for _, c := range m.palette.Colors() {
		hexStrings += hexColor(c) + "\n"

		if err != nil {
			return m, event.BuildDisplayCmd("error writing to palette file")
		}
	}

	_, err = f.WriteString(hexStrings)

	dir, _ := os.Getwd()
	msg := fmt.Sprintf("saved %s in /%s", filename, filepath.Base(dir))
	return m, event.BuildDisplayCmd(msg)
}

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
