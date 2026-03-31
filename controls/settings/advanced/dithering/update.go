package dithering

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

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
	Right: {
		DitherOn:         DitherOff,
		SerpentineOn:     SerpentineOff,
		ModeMatrix:       ModeBayer,
		ModeBayer:        ModeClusteredDot,
		MatrixList:       MatrixList,
		BayerSize2:       BayerSize4,
		BayerSize4:       BayerSize8,
		BayerSize8:       BayerSize16,
		ClusteredDotList: ClusteredDotList,
	},
	Left: {
		DitherOff:         DitherOn,
		SerpentineOff:     SerpentineOn,
		ModeClusteredDot:  ModeBayer,
		ModeBayer:         ModeMatrix,
		MatrixList:        MatrixList,
		BayerSize16:       BayerSize8,
		BayerSize8:        BayerSize4,
		BayerSize4:        BayerSize2,
		ClusteredDotList:  ClusteredDotList,
	},
	Up: {
		SerpentineOn:     DitherOn,
		SerpentineOff:    DitherOff,
		StrengthForm:     SerpentineOn,
		ModeMatrix:       StrengthForm,
		ModeBayer:        StrengthForm,
		ModeClusteredDot: StrengthForm,
		MatrixList:       ModeMatrix,
		BayerSize2:       ModeBayer,
		BayerSize4:       ModeBayer,
		BayerSize8:       ModeBayer,
		BayerSize16:      ModeBayer,
		ClusteredDotList: ModeClusteredDot,
	},
	Down: {
		DitherOn:         SerpentineOn,
		DitherOff:        SerpentineOff,
		SerpentineOn:     StrengthForm,
		SerpentineOff:    StrengthForm,
		StrengthForm:     ModeMatrix,
		ModeMatrix:       MatrixList,
		ModeBayer:        BayerSize2,
		ModeClusteredDot: ClusteredDotList,
	},
}

func (m Model) handleMatrixListUpdate(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Up) && m.matrixList.Index() == 0:
		return m.setFocus(ModeMatrix)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.setFocus(ModeMatrix)
	case key.Matches(msg, event.KeyMap.Enter):
		if it, ok := m.matrixList.SelectedItem().(matrixItem); ok {
			if mat, ok := errorDiffMatrixMap[it.Type]; ok {
				m.matrix = mat
			}
		}
		return m, event.StartRenderToViewCmd
	}
	var cmd tea.Cmd
	m.matrixList, cmd = m.matrixList.Update(msg)
	return m, cmd
}

func (m Model) handleClusteredDotListUpdate(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Up) && m.clusteredDotList.Index() == 0:
		return m.setFocus(ModeClusteredDot)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.setFocus(ModeClusteredDot)
	case key.Matches(msg, event.KeyMap.Enter):
		if it, ok := m.clusteredDotList.SelectedItem().(clusteredDotItem); ok {
			if cdm, ok := clusteredDotMatrixMap[it.Type]; ok {
				m.clusteredDot = cdm
			}
		}
		return m, event.StartRenderToViewCmd
	}
	var cmd tea.Cmd
	m.clusteredDotList, cmd = m.clusteredDotList.Update(msg)
	return m, cmd
}

func (m Model) handleStrengthUpdate(msg tea.Msg) (Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, event.KeyMap.Enter):
			m.strengthInput.Blur()
			m.active = m.focus
			return m, event.StartRenderToViewCmd
		case key.Matches(keyMsg, event.KeyMap.Esc):
			m.strengthInput.Blur()
			m.active = m.focus
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.strengthInput, cmd = m.strengthInput.Update(msg)
	return m, cmd
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	m.active = m.focus
	switch m.active {
	case DitherOn:
		m.doDithering = true
	case DitherOff:
		m.doDithering = false
	case SerpentineOn:
		m.doSerpentine = true
	case SerpentineOff:
		m.doSerpentine = false
	case ModeMatrix:
		m.ditherMode = Matrix
		m.modeControls = ModeMatrix
	case ModeBayer:
		m.ditherMode = Bayer
		m.modeControls = ModeBayer
	case ModeClusteredDot:
		m.ditherMode = ClusteredDot
		m.modeControls = ModeClusteredDot
	case BayerSize2:
		m.bayerSize = 2
	case BayerSize4:
		m.bayerSize = 4
	case BayerSize8:
		m.bayerSize = 8
	case BayerSize16:
		m.bayerSize = 16
	case StrengthForm:
		m.strengthInput.Focus()
		return m, nil
	default:
		// Items within a tab's content (list selections) are handled
		// by the list update handlers above
	}
	return m, event.StartRenderToViewCmd
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			return m.setFocus(next)
		}
	case key.Matches(msg, event.KeyMap.Up):
		if !m.doDithering && (m.focus == DitherOn || m.focus == DitherOff) {
			m.ShouldClose = true
			return m, nil
		}
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Down):
		if !m.doDithering {
			m.ShouldClose = true
			return m, nil
		}
		// Activate tab if focused but not active
		switch m.focus {
		case ModeMatrix:
			m.modeControls = ModeMatrix
			m.ditherMode = Matrix
		case ModeBayer:
			m.modeControls = ModeBayer
			m.ditherMode = Bayer
		case ModeClusteredDot:
			m.modeControls = ModeClusteredDot
			m.ditherMode = ClusteredDot
		}
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.ShouldClose = true
		}
	case key.Matches(msg, event.KeyMap.Tab):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			return m.setFocus(next)
		}
		if !m.doDithering {
			m.ShouldClose = true
			return m, nil
		}
		// Activate tab if focused but not active
		switch m.focus {
		case ModeMatrix:
			m.modeControls = ModeMatrix
			m.ditherMode = Matrix
		case ModeBayer:
			m.modeControls = ModeBayer
			m.ditherMode = Bayer
		case ModeClusteredDot:
			m.modeControls = ModeClusteredDot
			m.ditherMode = ClusteredDot
		}
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			return m.setFocus(next)
		} else {
			m.ShouldClose = true
		}
	}
	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Enter):
		return m.handleEnter()
	case key.Matches(msg, event.KeyMap.Nav):
		return m.handleNav(msg)
	case key.Matches(msg, event.KeyMap.Esc):
		return m.handleEsc()
	}
	return m, nil
}

func (m Model) setFocus(focus State) (Model, tea.Cmd) {
	m.focus = focus
	m.strengthInput.Blur()

	// Reset list cursors to the actual selection before deactivating
	m.resetMatrixListCursor()
	m.resetClusteredDotListCursor()

	m.matrixList.SetDelegate(NewDelegate(false))
	m.clusteredDotList.SetDelegate(NewDelegate(false))

	switch m.focus {
	case ModeMatrix, ModeBayer, ModeClusteredDot:
		// Tab focus only — don't change modeControls
	case MatrixList:
		m.matrixList.SetDelegate(NewDelegate(true))
	case ClusteredDotList:
		m.clusteredDotList.SetDelegate(NewDelegate(true))
	}
	return m, nil
}
