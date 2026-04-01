package app

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Zebbeni/ansizalizer/controls"
	"github.com/Zebbeni/ansizalizer/controls/modal"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/debug"
	"github.com/Zebbeni/ansizalizer/display"
	"github.com/Zebbeni/ansizalizer/env"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
	"github.com/Zebbeni/ansizalizer/viewer"
)

type State int

const (
	Splash State = iota
	Main
	Browser
	Settings
)

type Model struct {
	state        State
	splashTicks  int // total animation ticks elapsed (10ms each)
	splashSpawned int // how many characters have been spawned

	controls      controls.Model
	display       display.Model
	viewer        viewer.Model
	modal         *modal.Model // non-nil when a load/save modal is open
	quitConfirm   bool         // true when quit confirmation modal is showing
	quitFocusYes  bool         // true = Yes focused, false = No focused
	showHelp      bool         // true when help tooltips are visible

	waitingOnExport    bool
	exportQueue        []exportJob
	exportIndex        int
	overwriteModal     *OverwriteModal // non-nil when awaiting overwrite decision
	overwriteAllChoice OverwriteChoice // remembered choice when "Apply to all" is checked
	overwriteApproved  bool            // true when current file was approved for overwrite

	prefs prefs.Prefs

	w, h int
}

func New() Model {
	debug.Init()
	firstRun := prefs.Init()
	if firstRun {
		defaultCfg := settings.DefaultConfig()
		data, _ := json.MarshalIndent(defaultCfg, "", "  ")
		prefs.WriteDefaultSettings(data)
	}

	p := prefs.Load()

	if firstRun {
		p.App.ShowHelp = true
		p.Save()
	}

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

	startState := Splash
	if p.App.SkipSplash {
		startState = Main
	}

	m := Model{
		state:    startState,
		controls: controls.New(controlsWidth, p.Dirs),
		display:  display.New(),
		viewer:   viewer.New(),
		w:        100,
		h:        100,

		prefs: p,

		waitingOnExport: false,
	}

	// Restore settings from latest.json, default.json, or skip based on preferences
	if firstRun {
		m.controls.Settings.SaveLoad.SetStatus("Loaded default.json")
	} else if p.App.RestoreSettings {
		if cfg, err := settings.LoadConfig(prefs.LatestSettingsPath()); err == nil {
			m.controls.Settings.ApplyConfig(cfg)
			m.controls.Settings.SaveLoad.SetStatus("Loaded latest.json")
			debug.Log("Startup: loaded latest.json, chars mode=%s", cfg.Characters.Mode)
		}
	} else {
		if cfg, err := settings.LoadConfig(filepath.Join(prefs.SettingsDir(), "default.json")); err == nil {
			m.controls.Settings.ApplyConfig(cfg)
			m.controls.Settings.SaveLoad.SetStatus("Loaded default.json")
		}
	}

	// Sync palette to style package so paletted themes work at startup
	if m.controls.Settings.Colors.IsLimited() {
		style.PaletteColors = m.controls.Settings.Colors.GetCurrentPalette().Colors()
	}

	// Restore theme from preferences if enabled
	if p.App.RestoreTheme && p.App.LastTheme != "" {
		if themeName, ok := style.ThemeFromName[p.App.LastTheme]; ok {
			style.SetTheme(themeName)
			m.controls.Settings.Theme.SetThemeByName(p.App.LastTheme)
		}
	} else {
		style.SetTheme(style.ActiveTheme.Name)
	}

	// Apply show help preference
	m.showHelp = p.App.ShowHelp

	// Refresh all cached styles after theme is set
	m.controls.RefreshAllStyles()

	return m
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{splashTick()}
	// This initiates the polling cycle for window size updates
	// but shouldn't be necessary on non-Windows computers.
	if env.PollForSizeChange {
		cmds = append(cmds, pollForSizeChange)
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case splashTickMsg:
		if m.state != Splash {
			return m, nil
		}
		m.splashTicks++
		// Spawn a new character every spawnInterval ticks
		if m.splashSpawned < splashTotal && m.splashTicks%spawnInterval == 0 {
			m.splashSpawned++
		}
		// End once all characters have landed
		allSpawned := m.splashSpawned >= splashTotal
		ticksSinceLastSpawn := m.splashTicks - m.splashSpawned*spawnInterval
		if allSpawned && ticksSinceLastSpawn >= dropHeight {
			return m, splashPause()
		}
		return m, splashTick()
	case splashDoneMsg:
		m.state = Main
		return m, nil
	case checkSizeMsg:
		return m.handleCheckSizeMsg()
	case tea.WindowSizeMsg:
		return m.handleSizeMsg(msg)
	case event.CheckThemeMsg:
		return m.handleCheckTheme(msg)
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
		if m.state == Splash {
			m.state = Main
			return m, nil
		}
		// Quit confirmation modal captures all input
		if m.quitConfirm {
			return m.handleQuitModalKey(msg)
		}
		// Overwrite modal captures input during export
		if m.overwriteModal != nil {
			m.overwriteModal.Update(msg)
			if m.overwriteModal.choice != OverwriteUndecided {
				return m.handleOverwriteChoice()
			}
			return m, nil
		}
		// Load/save modal captures all input when open
		if m.modal != nil {
			return m.handleModalUpdate(msg)
		}
		// Ctrl+Q opens quit confirmation from anywhere
		if key.Matches(msg, event.KeyMap.Quit) {
			m.quitConfirm = true
			m.quitFocusYes = true
			return m, nil
		}
		// Ctrl+H toggles help tooltips
		if key.Matches(msg, event.KeyMap.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}
		switch {
		case key.Matches(msg, event.KeyMap.Copy):
			return m.handleCopy()
		case key.Matches(msg, event.KeyMap.Save):
			return m.handleSave()
		case key.Matches(msg, event.KeyMap.Debug):
			return m.handleDebug()
		}
	}
	if m.state == Splash {
		return m, nil
	}
	// Modal also captures non-key messages
	if m.modal != nil {
		return m.handleModalUpdate(msg)
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
func (m Model) View() tea.View {
	if m.state == Splash {
		sv := tea.NewView(m.renderSplash())
		sv.AltScreen = true
		return sv
	}

	controls := style.ApplyBg(m.renderControls(), 0)
	display := style.ApplyBg(m.display.View(), 0)
	viewer := style.ApplyBg(m.renderViewer(), 0)
	help := m.renderHelp() // already has ApplyBg baked into cache

	leftPanel := controls
	rightPanel := lipgloss.JoinVertical(lipgloss.Top, display, viewer)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	all := lipgloss.JoinVertical(lipgloss.Top, panels, help)

	appStyle := lipgloss.NewStyle().Width(m.w).Height(m.h)
	if !style.ActiveTheme.Transparent {
		appStyle = appStyle.Background(style.ActiveTheme.Bg)
	}

	rendered := appStyle.Render(all)

	// Overlay modal centered on screen using Canvas
	if m.modal != nil && m.modal.Open {
		modalView := m.modal.View()
		modalW := lipgloss.Width(modalView)
		modalH := lipgloss.Height(modalView)
		mx := (m.w - modalW) / 2
		my := (m.h - modalH) / 2

		canvas := lipgloss.NewCanvas(m.w, m.h)
		base := lipgloss.NewLayer(rendered)
		overlay := lipgloss.NewLayer(modalView).X(mx).Y(my).Z(1)
		canvas.Compose(lipgloss.NewCompositor(base, overlay))
		rendered = canvas.Render()
	}

	// Overlay file dropdown using lipgloss Canvas if open
	if m.controls.FileDropdown.Open {
		dropY := m.controls.ButtonHeight()

		canvas := lipgloss.NewCanvas(m.w, m.h)
		base := lipgloss.NewLayer(rendered)
		mainMenu := lipgloss.NewLayer(m.controls.FileDropdown.MainMenuView()).X(1).Y(dropY).Z(1)

		layers := []*lipgloss.Layer{base, mainMenu}

		if sub := m.controls.FileDropdown.SubMenuView(); sub != "" {
			subX := 1 + m.controls.FileDropdown.MainMenuWidth()
			subY := dropY + m.controls.FileDropdown.SubMenuYOffset()
			subMenu := lipgloss.NewLayer(sub).X(subX).Y(subY).Z(2)
			layers = append(layers, subMenu)
		}

		canvas.Compose(lipgloss.NewCompositor(layers...))
		rendered = canvas.Render()
	}

	// Overlay help tooltip to the right of the controls panel
	if m.showHelp && !m.controls.FileDropdown.Open {
		if helpText := m.controls.HelpText(); helpText != "" {
			tooltip := style.RenderTooltip(helpText)
			tx := controlsWidth + 1
			ty := m.controls.ButtonHeight() + 1

			canvas := lipgloss.NewCanvas(m.w, m.h)
			base := lipgloss.NewLayer(rendered)
			overlay := lipgloss.NewLayer(tooltip).X(tx).Y(ty).Z(5)
			canvas.Compose(lipgloss.NewCompositor(base, overlay))
			rendered = canvas.Render()
		}
	}

	// Overlay overwrite confirmation modal during export
	if m.overwriteModal != nil {
		owView := m.overwriteModal.View()
		owW := lipgloss.Width(owView)
		owH := lipgloss.Height(owView)
		owX := (m.w - owW) / 2
		owY := (m.h - owH) / 2

		canvas := lipgloss.NewCanvas(m.w, m.h)
		base := lipgloss.NewLayer(rendered)
		overlay := lipgloss.NewLayer(owView).X(owX).Y(owY).Z(8)
		canvas.Compose(lipgloss.NewCompositor(base, overlay))
		rendered = canvas.Render()
	}

	// Overlay quit confirmation modal
	if m.quitConfirm {
		quitView := m.renderQuitModal()
		qw := lipgloss.Width(quitView)
		qh := lipgloss.Height(quitView)
		qx := (m.w - qw) / 2
		qy := (m.h - qh) / 2

		canvas := lipgloss.NewCanvas(m.w, m.h)
		base := lipgloss.NewLayer(rendered)
		overlay := lipgloss.NewLayer(quitView).X(qx).Y(qy).Z(10)
		canvas.Compose(lipgloss.NewCompositor(base, overlay))
		rendered = canvas.Render()
	}

	v := tea.NewView(rendered)
	v.AltScreen = true
	return v
}
