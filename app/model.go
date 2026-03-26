package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Zebbeni/ansizalizer/controls"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/display"
	"github.com/Zebbeni/ansizalizer/env"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
	"github.com/Zebbeni/ansizalizer/viewer"
)

type State int

const (
	Main State = iota
	Browser
	Settings
)

type Model struct {
	state State

	controls controls.Model
	display  display.Model
	viewer   viewer.Model

	waitingOnExport bool
	exportQueue     []exportJob
	exportIndex     int

	prefs prefs.Prefs

	w, h int
}

func New() Model {
	firstRun := prefs.Init()
	if firstRun {
		defaultCfg := settings.DefaultConfig()
		data, _ := json.MarshalIndent(defaultCfg, "", "  ")
		prefs.WriteDefaultSettings(data)
	}

	p := prefs.Load()

	// Set default dirs
	if p.Dirs.Browse == "" {
		if home, err := os.UserHomeDir(); err == nil {
			docs := filepath.Join(home, "Documents")
			if info, err := os.Stat(docs); err == nil && info.IsDir() {
				p.Dirs.Browse = docs
			}
		}
	}
	if p.Dirs.SaveDir == "" {
		p.Dirs.SaveDir = prefs.SettingsDir()
	}
	if p.Dirs.LoadFile == "" {
		p.Dirs.LoadFile = prefs.SettingsDir()
	}
	if p.Dirs.PaletteLoad == "" {
		p.Dirs.PaletteLoad = prefs.PalettesDir()
	}
	p.Save()

	m := Model{
		state:    Main,
		controls: controls.New(controlsWidth, p.Dirs),
		display:  display.New(),
		viewer:   viewer.New(),
		w:        100,
		h:        100,

		prefs: p,

		waitingOnExport: false,
	}

	// Restore settings from latest.json (or default.json on first run)
	if firstRun {
		m.controls.Settings.SaveLoad.SetStatus("Loaded default.json")
	} else if cfg, err := settings.LoadConfig(prefs.LatestSettingsPath()); err == nil {
		m.controls.Settings.ApplyConfig(cfg)
		m.controls.Settings.SaveLoad.SetStatus("Loaded latest.json")
	}

	// Sync palette to style package so paletted themes work at startup
	if m.controls.Settings.Colors.IsLimited() {
		style.PaletteColors = m.controls.Settings.Colors.GetCurrentPalette().Colors()
	}
	style.SetTheme(style.ActiveTheme.Name)

	return m
}

func (m Model) Init() tea.Cmd {
	// This initiates the polling cycle for window size updates
	// but shouldn't be necessary on non-Windows computers.
	if env.PollForSizeChange {
		return pollForSizeChange
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkSizeMsg:
		return m.handleCheckSizeMsg()
	case tea.WindowSizeMsg:
		return m.handleSizeMsg(msg)
	case event.StartRenderToViewMsg:
		return m.handleStartRenderToViewCmd()
	case event.FinishRenderToViewMsg:
		return m.handleFinishRenderToViewMsg(msg)
	case event.FinishRenderGIFToViewMsg:
		return m.handleFinishRenderGIFToViewMsg(msg)
	case event.AnimationTickMsg:
		var cmd tea.Cmd
		m.viewer, cmd = m.viewer.Update(msg)
		return m, cmd
	case event.StartAdaptingMsg:
		return m.handleStartAdaptingMsg()
	case event.FinishAdaptingMsg:
		return m.handleFinishAdaptingMsg(msg)
	case event.LospecRequestMsg:
		return m.handleLospecRequestMsg(msg)
	case event.LospecResponseMsg:
		return m.handleLospecResponseMsg(msg)
	case event.StartExportMsg:
		return m.handleStartExportMsg(msg)
	case event.StartRenderToExportMsg:
		return m.handleRenderToExportMsg()
	case event.DisplayMsg:
		return m.handleDisplayMsg(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, event.KeyMap.Copy):
			return m.handleCopy()
		case key.Matches(msg, event.KeyMap.Save):
			return m.handleSave()
		case key.Matches(msg, event.KeyMap.Debug):
			return m.handleDebug()
		}
	}
	return m.handleControlsUpdate(msg)
}

// View puts the whole TUI together, laid out like this:
//
//	(Left Panel)                (Right Panel)
//
// ┌────────────────┬────────────────────────────────────────┐
// │   Controls     │               Display                  │
// │                ├────────────────────────────────────────┤
// │                │               Viewer                   │
// │                │                                        │
// ├────────────────┴────────────────────────────────────────┤
// │               Help                                      │
// └─────────────────────────────────────────────────────────┘
func (m Model) View() string {
	controls := style.ApplyBg(m.renderControls(), 0)
	display := style.ApplyBg(m.display.View(), 0)
	viewer := m.renderViewer()
	help := m.renderHelp() // already has ApplyBg baked into cache

	leftPanel := controls
	rightPanel := lipgloss.JoinVertical(lipgloss.Top, display, viewer)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	all := lipgloss.JoinVertical(lipgloss.Top, panels, help)

	appStyle := lipgloss.NewStyle().Width(m.w).Height(m.h)
	if !style.ActiveTheme.Transparent {
		appStyle = appStyle.Background(style.ActiveTheme.Bg)
	}

	return appStyle.Render(all)
}
