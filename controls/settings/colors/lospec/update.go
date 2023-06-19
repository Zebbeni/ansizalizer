package lospec

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
)

var (
	navMap = map[Direction]map[State]State{
		Right: {CountForm: TagForm, FilterAny: FilterExact, FilterExact: FilterMax, FilterMax: FilterMin, SortAlphabetical: SortDownloads, SortDownloads: SortNewest},
		Left:  {TagForm: CountForm, FilterMin: FilterMax, FilterMax: FilterExact, FilterExact: FilterAny, SortNewest: SortDownloads, SortDownloads: SortAlphabetical},
		Up:    {FilterAny: CountForm, FilterExact: CountForm, FilterMax: TagForm, FilterMin: TagForm, SortAlphabetical: FilterAny, SortDownloads: TagForm, SortNewest: TagForm, List: SortAlphabetical},
		Down:  {CountForm: FilterAny, TagForm: FilterMax, FilterAny: SortAlphabetical, FilterExact: SortDownloads, FilterMax: SortDownloads, FilterMin: SortNewest, SortAlphabetical: List, SortDownloads: List, SortNewest: List},
	}
	filterParams = map[State]string{
		FilterAny:   "any",
		FilterExact: "exact",
		FilterMax:   "max",
		FilterMin:   "min",
	}
	sortParams = map[State]string{
		SortAlphabetical: "alphabetical",
		SortDownloads:    "downloads",
		SortNewest:       "newest",
	}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	switch m.focus {
	case CountForm:
		m.countInput.Focus()
		return m, nil
	case TagForm:
		m.tagInput.Focus()
		return m, nil
	case FilterAny, FilterExact, FilterMax, FilterMin:
		m.filterType = m.focus
		return m.searchLospec(0)
	case SortAlphabetical, SortDownloads, SortNewest:
		m.sortType = m.focus
		return m.searchLospec(0)
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
		m.palettes = make([]list.Item, msg.Data.TotalCount)
		m.paletteList = CreateList(m.palettes, m.width-2)
		m.isPaletteListAllocated = true
	}

	// use the page number*10 (assumes 10 palettes per page) to populate palettes
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
		m.palettes[idx] = palette.New(p.Title, colors, m.width-4, 2)
	}

	m.paletteList.SetItems(m.palettes)

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

func (m Model) handleListUpdate(msg tea.Msg) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	index := m.paletteList.Index()
	if key.Matches(keyMsg, event.KeyMap.Up) && index == 0 {
		return m.handleNav(keyMsg)
	}

	var cmd tea.Cmd
	m.paletteList, cmd = m.paletteList.Update(msg)
	return m, cmd
}

func (m Model) searchLospec(page int) (Model, tea.Cmd) {
	colors, _ := strconv.Atoi(m.countInput.Value())
	tag := m.tagInput.Value()
	filterType := filterParams[m.filterType]
	sortingType := "alphabetical"
	//sortingType := sortParams[m.sortType]
	urlString := "https://lospec.com/palette-list/load?colorNumber=%d&tag=%s&colorNumberFilterType=%s&sortingType=%s&page=%d"
	url := fmt.Sprintf(urlString, colors, tag, filterType, sortingType, page)
	return m, event.BuildLospecRequestCmd(event.LospecRequestMsg{
		URL:  url,
		ID:   m.requestID,
		Page: page,
	})
}
