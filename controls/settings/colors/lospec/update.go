package lospec

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/palette"
)

// TODO: Direction is redefined in multiple places

type Direction int

type Param int

const (
	Left Direction = iota
	Right
	Up
	Down

	Colors Param = iota
	Tag
	Filter
	Sorting
	//lospecRequestUrl = "https://lospec.com/palette-list/load?colorNumberFilterType=exact&colorNumber=32&tag=&sortingType=alphabetical"
	lospecRequestUrl = "https://lospec.com/palette-list/load?colorNumberFilterType=exact&colorNumber=32&tag=&sortingType=alphabetical"
)

var (
	navMap = map[Direction]map[State]State{
		Right: {CountForm: TagForm},
		Left:  {TagForm: CountForm},
		Up:    {},
		Down:  {},
	}
)

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
	case TagForm:
		m.tagInput.Focus()
		return m, nil
	}
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
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
	return m, nil
}

func (m Model) handleLospecResponse(msg event.LospecResponseMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	// return early if response no longer matches current requestID
	if msg.ID != m.requestID {
		return m, cmd
	}

	// if we haven't initialized and allocated an array of palettes for the current , do that first
	if !m.isPaletteListAllocated {
		m.paletteList = make([]palette.Model, msg.Data.TotalCount)
		m.isPaletteListAllocated = true
	}

	// use the page number (and assuming 10 palettes per page) to populate our paletteList
	for i, p := range msg.Data.Palettes {
		colors := make([]color.Color, len(p.Colors))
		var err error

		for colorIndex, c := range p.Colors {
			colors[colorIndex], err = colorful.Hex(fmt.Sprintf("#%s", c))
			if err != nil {
				return m, event.BuildDisplayCmd("error converting hex value")
			}

		}

		idx := (msg.Page * 10) + i
		m.paletteList[idx] = palette.New(p.Title, colors, m.width-2, 1)
	}

	return m, cmd
}

func (m Model) handleCountFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.countInput, cmd = m.countInput.Update(msg)

	// If Update caused countInput to become de-focused, kick off a new request
	if m.countInput.Focused() == false {
		m.requestID += 1
		m.isPaletteListAllocated = false
		return m.searchLospec(0)
	}
	return m, cmd
}

func (m Model) handleTagFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.tagInput, cmd = m.tagInput.Update(msg)
	if m.tagInput.Focused() == false {
		m.requestID += 1
		m.isPaletteListAllocated = false
		return m.searchLospec(0)
	}
	return m, cmd
}

func (m Model) searchLospec(page int) (Model, tea.Cmd) {
	colors, _ := strconv.Atoi(m.countInput.Value())
	tag := m.tagInput.Value()
	filterType := "exact"
	sortingType := "alphabetical"
	urlString := "https://lospec.com/palette-list/load?colorNumber=%d&tag=%s&colorNumberFilterType=%s&sortingType=%s&page=%d"
	url := fmt.Sprintf(urlString, colors, tag, filterType, sortingType, page)
	return m, event.BuildLospecRequestCmd(event.LospecRequestMsg{
		URL:  url,
		ID:   m.requestID,
		Page: page,
	})
}
