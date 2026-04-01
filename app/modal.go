package app

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"

	"github.com/Zebbeni/ansizalizer/controls/modal"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) openLoadModal() Model {
	md := modal.NewLoad(prefs.SettingsDir())
	m.modal = &md
	return m
}

func (m Model) openSaveModal() Model {
	md := modal.NewSave(prefs.SettingsDir())
	m.modal = &md
	return m
}

func (m Model) openExportFileModal() Model {
	sourceFile := m.controls.FileBrowser.SelectedFile
	destDir := m.controls.Export.Destination.GetSelected()
	md := modal.NewExportFile(sourceFile, destDir)
	m.modal = &md
	return m
}

func (m Model) openPrefsModal() Model {
	md := modal.NewPrefs(m.prefs.App.ShowHelp, m.prefs.App.RestoreTheme, m.prefs.App.RestoreSettings)
	m.modal = &md
	return m
}

func (m Model) openExportBatchModal() Model {
	sourceDir := m.controls.FileBrowser.SelectedDir
	destDir := m.controls.Export.Destination.GetSelected()
	md := modal.NewExportBatch(sourceDir, destDir)
	m.modal = &md
	return m
}

func (m Model) handleModalUpdate(msg tea.Msg) (Model, tea.Cmd) {
	md, cmd := m.modal.Update(msg)
	m.modal = &md

	// Check if modal produced a result
	if result := m.modal.GetResult(); result != nil {
		return m.handleModalResult(*result)
	}

	// Check if modal was closed (cancelled)
	if !m.modal.Open {
		m.modal = nil
	}

	return m, cmd
}

func (m Model) handleModalResult(result modal.Result) (Model, tea.Cmd) {
	m.modal = nil

	switch result.Kind {
	case modal.LoadKind:
		cfg, err := settings.LoadConfig(result.Path)
		if err != nil {
			return m, event.BuildDisplayCmd("error loading settings")
		}
		m.controls.Settings.ApplyConfig(cfg)
		filename := filepath.Base(result.Path)

		// Sync palette colors to the global style system so paletted themes work
		cmds := []tea.Cmd{event.StartRenderToViewCmd, event.BuildDisplayCmd("loaded " + filename)}
		if m.controls.Settings.Colors.IsLimited() {
			p := m.controls.Settings.Colors.GetCurrentPalette().Colors()
			if len(p) > 0 {
				style.PaletteColors = p
			}
		}
		// Refresh the current theme with the new palette
		cmds = append(cmds, event.BuildCheckThemeCmd(style.PaletteColors, style.ThemeNames[style.ActiveTheme.Name]))
		return m, tea.Batch(cmds...)

	case modal.SaveKind:
		cfg := m.controls.Settings.ExportConfig()
		err := settings.SaveConfig(cfg, result.Path)
		if err != nil {
			return m, event.BuildDisplayCmd("error saving settings")
		}
		filename := filepath.Base(result.Path)
		return m, event.BuildDisplayCmd("saved " + filename)

	case modal.PrefsKind:
		m.prefs.App.ShowHelp = result.ShowHelp
		m.prefs.App.RestoreTheme = result.RestoreTheme
		m.prefs.App.RestoreSettings = result.RestoreSettings
		m.showHelp = result.ShowHelp
		if result.RestoreTheme {
			m.prefs.App.LastTheme = style.ThemeNames[style.ActiveTheme.Name]
		}
		m.prefs.Save()
		return m, event.BuildDisplayCmd("preferences saved")

	case modal.ExportFileKind, modal.ExportBatchKind:
		return m, event.BuildStartExportCmd(event.StartExportMsg{
			SourcePath:      result.SourcePath,
			DestinationPath: result.DestinationPath,
			IsDir:           result.IsDir,
			UseSubDirs:      result.UseSubDirs,
		})
	}

	return m, nil
}
