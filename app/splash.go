package app

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	splashRow1 = []rune{' ', '▛', '▜', '▐', '▜', '▐', '▀', '▐', '▝', '▜', '▐', '▜', '▐', ' ', '▐', '▝', '▜', '▐', '▛', '▐', '▀', '▌'}
	splashRow2 = []rune{' ', '▛', '▜', '▐', '▐', '▗', '▟', '▐', '▐', '▄', '▐', '▜', '▐', '▄', '▐', '▐', '▄', '▐', '▄', '▐', '▀', '▖'}

	splashCols = len(splashRow1)
	// Total characters to spawn (bottom-top order: row2[0], row1[0], row2[1], row1[1], ...)
	splashTotal   = splashCols * 2
	dropHeight    = 2
	spawnInterval = 5 // ticks between spawns (50ms at 10ms/tick)
)

type splashTickMsg struct{}
type splashDoneMsg struct{}

func splashTick() tea.Cmd {
	return tea.Tick(5*time.Millisecond, func(t time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

func splashPause() tea.Cmd {
	return tea.Tick(1000*time.Millisecond, func(t time.Time) tea.Msg {
		return splashDoneMsg{}
	})
}

func (m Model) renderSplash() string {
	// Grid: rows 0-1 are the drop zone above, rows 2-3 are the final title positions.
	// splashRow1 lands at grid row 2, splashRow2 lands at grid row 3.
	gridRows := 2 + dropHeight
	grid := make([][]rune, gridRows)
	for r := range grid {
		grid[r] = make([]rune, splashCols)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	// Place each spawned character at its current falling position.
	// Character i was spawned at tick i*spawnInterval. Elapsed ticks since
	// spawn = splashTicks - i*spawnInterval. It drops 1 row per tick.
	for i := 0; i < m.splashSpawned; i++ {
		col := i / 2
		var finalRow int
		var ch rune
		if i%2 == 0 {
			finalRow = 3
			ch = splashRow2[col]
		} else {
			finalRow = 2
			ch = splashRow1[col]
		}
		ticksSinceSpawn := m.splashTicks - (i+1)*spawnInterval
		offset := dropHeight - ticksSinceSpawn
		if offset < 0 {
			offset = 0
		}
		row := finalRow - offset
		if row >= 0 && row < gridRows {
			grid[row][col] = ch
		}
	}

	// Render the grid rows and copyright, all at consistent width
	copyrightText := "\u00a9 2026 Simulacrum Studios, LLC"
	titleStyle := style.BgStyle().Foreground(style.NormalColor1)
	copyrightStyle := style.BgStyle().Foreground(style.DimmedColor1)

	// Use the wider of title or copyright as the block width
	blockWidth := lipgloss.Width(titleStyle.Render(string(grid[0])))
	if cw := lipgloss.Width(copyrightStyle.Render(copyrightText)); cw > blockWidth {
		blockWidth = cw
	}

	rows := make([]string, gridRows)
	for r := 0; r < gridRows; r++ {
		rows[r] = titleStyle.Width(blockWidth).AlignHorizontal(lipgloss.Center).Render(string(grid[r]))
	}
	title := lipgloss.JoinVertical(lipgloss.Left, rows...)

	var copyrightRow string
	if m.splashSpawned >= splashTotal {
		copyrightRow = copyrightStyle.Width(blockWidth).AlignHorizontal(lipgloss.Center).PaddingTop(1).Render(copyrightText)
	} else {
		copyrightRow = style.BgStyle().Width(blockWidth).Render("")
	}
	splash := lipgloss.JoinVertical(lipgloss.Left, title, copyrightRow)

	appStyle := lipgloss.NewStyle().Width(m.w).Height(m.h)
	if !style.ActiveTheme.Transparent {
		appStyle = appStyle.Background(style.ActiveTheme.Bg)
	}

	centered := lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, splash,
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(style.ActiveTheme.Bg)))

	return appStyle.Render(centered)
}
