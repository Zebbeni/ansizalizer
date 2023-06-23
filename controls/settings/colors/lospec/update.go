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
	"github.com/Zebbeni/ansizalizer/style"
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
		Right: {CountForm: FilterExact, FilterExact: FilterMax, FilterMax: FilterMin, SortAlphabetical: SortDownloads, SortDownloads: SortNewest},
		Left:  {TagForm: CountForm, FilterMin: FilterMax, FilterMax: FilterExact, FilterExact: CountForm, SortNewest: SortDownloads, SortDownloads: SortAlphabetical},
		Up:    {TagForm: CountForm, SortAlphabetical: TagForm, SortDownloads: TagForm, SortNewest: TagForm, List: SortAlphabetical},
		Down:  {CountForm: TagForm, FilterExact: TagForm, FilterMax: TagForm, FilterMin: TagForm, TagForm: SortAlphabetical, SortAlphabetical: List, SortDownloads: List, SortNewest: List},
	}
	filterParams = map[State]string{
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
	case FilterExact, FilterMax, FilterMin:
		m.filterType = m.focus
		return m.searchLospec(0)
	case SortAlphabetical, SortDownloads, SortNewest:
		m.sortType = m.focus
		return m.searchLospec(0)
	case List:
		m.palette, _ = m.paletteList.SelectedItem().(palette.Model)
		return m, event.StartRenderCmd
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

	// if we haven't initialized and allocated an array of palettes for the current request series, do that first
	if !m.isPaletteListAllocated {
		m.palettes = make([]list.Item, msg.Data.TotalCount)
		m.paletteList = CreateList(m.palettes, m.width-2)
		m.paletteList.Styles.Title = style.DimmedTitle.Padding(0, 1, 0, 2).Width(25).MaxWidth(25)
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
		return m.searchLospec(0)
	}
	return m, cmd
}

func (m Model) handleTagFormUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.tagInput, cmd = m.tagInput.Update(msg)
	if m.tagInput.Focused() == false {
		return m.searchLospec(0)
	}
	return m, cmd
}

func (m Model) handleListUpdate(msg tea.Msg) (Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		return m.handleEnter()
	case key.Matches(keyMsg, event.KeyMap.Up) && m.paletteList.Index() == 0:
		return m.handleNav(keyMsg)
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.focus = TagForm
	}

	var cmd tea.Cmd
	m.paletteList, cmd = m.paletteList.Update(msg)
	if m.paletteList.Index() < (m.highestPageRequested-1)*10 {
		return m, cmd
	}

	m.highestPageRequested += 1
	return m.searchLospec(m.highestPageRequested)
}

func (m Model) searchLospec(page int) (Model, tea.Cmd) {
	if page == 0 {
		m.requestID += 1
		m.highestPageRequested = 0
		m.isPaletteListAllocated = false
	}

	colors, _ := strconv.Atoi(m.countInput.Value())
	tag := m.tagInput.Value()
	filterType := filterParams[m.filterType]
	sortingType := "alphabetical"

	urlString := "https://lospec.com/palette-list/load?colorNumber=%d&tag=%s&colorNumberFilterType=%s&sortingType=%s&page=%d"
	url := fmt.Sprintf(urlString, colors, tag, filterType, sortingType, page)
	return m, event.BuildLospecRequestCmd(event.LospecRequestMsg{
		URL:  url,
		ID:   m.requestID,
		Page: page,
	})
}
