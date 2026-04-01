package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	tea "charm.land/bubbletea/v2"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/app/adapt"
	"github.com/Zebbeni/ansizalizer/app/process"
	"github.com/Zebbeni/ansizalizer/controls/filemenu"
	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/debug"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/prefs"
	"github.com/Zebbeni/ansizalizer/style"
)

func (m Model) handleCheckTheme(msg event.CheckThemeMsg) (Model, tea.Cmd) {
	prevBg := style.ColorHex(style.ActiveTheme.Bg)
	prevTheme := style.ActiveTheme.Name

	// Update palette colors (don't clear if empty — may be temporarily empty)
	if len(msg.Palette) > 0 {
		style.PaletteColors = msg.Palette
	}

	// Resolve and apply theme
	if msg.ThemeName != "" {
		if themeName, ok := style.ThemeFromName[msg.ThemeName]; ok {
			style.SetTheme(themeName)
			// Update the theme submenu to match
			m.controls.Settings.Theme.SetThemeByName(msg.ThemeName)
		}
	} else if prevTheme == style.LightOnDarkPaletted || prevTheme == style.DarkOnLightPaletted {
		// Refresh paletted theme with new palette
		style.SetTheme(prevTheme)
	}

	newBg := style.ColorHex(style.ActiveTheme.Bg)
	var cmds []tea.Cmd

	if prevBg != newBg || prevTheme != style.ActiveTheme.Name {
		debug.Log("handleCheckTheme: theme changed (bg %s -> %s)", prevBg, newBg)
		m.controls.RefreshAllStyles()
		helpCacheBg = ""
		cmds = append(cmds, event.StartRenderToViewCmd)

		// Persist last theme for restore-on-start
		if m.prefs.App.RestoreTheme {
			m.prefs.App.LastTheme = style.ThemeNames[style.ActiveTheme.Name]
			m.prefs.Save()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleStartRenderToViewCmd() (Model, tea.Cmd) {
	m.viewer.WaitingOnRender = true
	m.saveLatestSettings()
	return m, m.processRenderToViewCmd
}

func (m Model) handleFinishRenderToViewMsg(msg event.FinishRenderToViewMsg) (Model, tea.Cmd) {
	// cut out early if the finished render is for a previously selected image
	if msg.FilePath != m.controls.FileBrowser.SelectedFile {
		return m, nil
	}
	var cmd tea.Cmd
	m.controls.Settings.Alpha.AlphaImage = false
	m.controls.Settings.Animation.AnimatedImage = false
	re := regexp.MustCompile(`(?i)\.(png|gif)$`)
	if re.Match([]byte(m.controls.FileBrowser.SelectedFile)) {
		m.controls.Settings.Alpha.AlphaImage = true
	}
	m.viewer, cmd = m.viewer.Update(msg)
	return m, cmd
}

func (m Model) processRenderToViewCmd() tea.Msg {
	re := regexp.MustCompile(`(?i)\.(png|gif)$`)
	filePath := m.controls.FileBrowser.SelectedFile

	colorsString := "true color"
	alphaString := ""
	if m.controls.Settings.Colors.IsLimited() {
		palette := m.controls.Settings.Colors.GetCurrentPalette()
		colorsString = palette.Title()
	}
	if m.controls.Settings.Alpha.ShouldOutputAlpha() && re.Match([]byte(filePath)) {
		alphaString = "alpha channel"
	}

	brightness := m.controls.Settings.Adjust.Brightness()
	contrast := m.controls.Settings.Adjust.Contrast()
	adjustParts := make([]string, 0, 2)
	if brightness != 0 {
		adjustParts = append(adjustParts, fmt.Sprintf("bright %d", brightness))
	}
	if contrast != 0 {
		adjustParts = append(adjustParts, fmt.Sprintf("contrast %d", contrast))
	}
	adjustString := ""
	if len(adjustParts) > 0 {
		adjustString = strings.Join(adjustParts, ", ")
	}

	doDither, _, _ := m.controls.Settings.Advanced.Dithering()
	ditherString := ""
	if doDither {
		ditherString = "dithered"
	}

	extraParts := make([]string, 0, 2)
	if adjustString != "" {
		extraParts = append(extraParts, adjustString)
	}
	if ditherString != "" {
		extraParts = append(extraParts, ditherString)
	}
	extraInfo := strings.Join(extraParts, ", ")

	// Compute solid bg for non-transparent themes
	var solidBg *colorful.Color
	if !style.ActiveTheme.Transparent {
		hex := style.ColorHex(style.ActiveTheme.Bg)
		c, err := colorful.Hex(hex)
		if err == nil {
			solidBg = &c
		}
	}

	// Animated GIF path
	if process.IsGIFFile(filePath) {
		frames, _ := process.RenderGIFFileWithBg(m.controls.Settings, filePath, solidBg)
		if len(frames) > 1 {
			return event.FinishRenderGIFToViewMsg{
				FilePath:     filePath,
				Frames:       frames,
				Delay:        m.controls.Settings.Animation.Delay(),
				ColorsString: colorsString,
				AlphaString:  alphaString,
				ExtraInfo:    extraInfo,
			}
		}
		// Single-frame GIF: use normal path
		return event.FinishRenderToViewMsg{
			FilePath:     filePath,
			ImgString:    frames[0],
			ColorsString: colorsString,
			AlphaString:  alphaString,
			ExtraInfo:    extraInfo,
		}
	}

	// Non-GIF: existing path
	imgString := process.RenderImageFileWithBg(m.controls.Settings, filePath, solidBg)
	return event.FinishRenderToViewMsg{
		FilePath:     filePath,
		ImgString:    imgString,
		ColorsString: colorsString,
		AlphaString:  alphaString,
		ExtraInfo:    extraInfo,
	}
}

func (m Model) handleFinishRenderGIFToViewMsg(msg event.FinishRenderGIFToViewMsg) (Model, tea.Cmd) {
	if msg.FilePath != m.controls.FileBrowser.SelectedFile {
		return m, nil
	}
	m.controls.Settings.Alpha.AlphaImage = true
	m.controls.Settings.Animation.AnimatedImage = true
	var cmd tea.Cmd
	m.viewer, cmd = m.viewer.Update(msg)
	return m, cmd
}

func (m Model) handleStartExportMsg(msg event.StartExportMsg) (Model, tea.Cmd) {
	if m.waitingOnExport {
		return m, nil
	}

	var exportQueue []exportJob
	var err error

	// build export queue
	if msg.IsDir {
		exportQueue, err = buildExportQueue(msg.SourcePath, msg.DestinationPath, msg.UseSubDirs)
		if err != nil {
			return m, event.BuildDisplayCmd(fmt.Sprintf("error exporting: %s", err))
		}
	} else {
		nameWithoutExt := strings.Split(filepath.Base(msg.SourcePath), ".")[0]
		destFilePath := filepath.Join(msg.DestinationPath, fmt.Sprintf("%s.ansi", nameWithoutExt))
		exportQueue = []exportJob{
			{
				sourcePath:      msg.SourcePath,
				destinationPath: destFilePath,
			},
		}
	}

	m.exportIndex = 0
	m.exportQueue = exportQueue
	m.waitingOnExport = true
	m.overwriteAllChoice = OverwriteUndecided
	m.overwriteModal = nil

	return m, tea.Batch(event.StartRenderToExportCmd, event.BuildDisplayCmd(fmt.Sprintf("%d jobs queued", len(exportQueue))))
}

func (m Model) skipCurrentExport() (Model, tea.Cmd) {
	m.exportIndex++
	displayMsg := fmt.Sprintf("%d/%d exports completed (skipped)", m.exportIndex, len(m.exportQueue))
	displayCmd := event.BuildDisplayCmd(displayMsg)

	if m.exportIndex >= len(m.exportQueue) {
		m.waitingOnExport = false
		return m, displayCmd
	}
	return m, tea.Batch(event.StartRenderToExportCmd, displayCmd)
}

func (m Model) handleRenderToExportMsg() (Model, tea.Cmd) {

	currentJob := m.exportQueue[m.exportIndex]

	// Check if destination already exists
	destToCheck := currentJob.destinationPath
	if process.IsGIFFile(currentJob.sourcePath) {
		destToCheck = strings.TrimSuffix(currentJob.destinationPath, ".ansi")
	}
	if _, err := os.Stat(destToCheck); err == nil && !m.overwriteApproved {
		if m.overwriteAllChoice == OverwriteSkip {
			return m.skipCurrentExport()
		} else if m.overwriteAllChoice != OverwriteYes {
			m.overwriteModal = newOverwriteModal(destToCheck)
			return m, nil
		}
	}

	if process.IsGIFFile(currentJob.sourcePath) {
		frames, _ := process.RenderGIFFile(m.controls.Settings, currentJob.sourcePath)
		if len(frames) > 1 {
			// create directory and save each frame
			destDir := strings.TrimSuffix(currentJob.destinationPath, ".ansi")
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return m, event.BuildDisplayCmd("error creating export directory")
			}
			baseName := strings.TrimSuffix(filepath.Base(currentJob.destinationPath), ".ansi")
			for i, frame := range frames {
				framePath := filepath.Join(destDir, fmt.Sprintf("%s_%03d.ansi", baseName, i+1))
				if err := writeFile(framePath, frame); err != nil {
					return m, event.BuildDisplayCmd(fmt.Sprintf("error writing frame %d", i+1))
				}
			}
		} else {
			if err := writeFile(currentJob.destinationPath, frames[0]); err != nil {
				return m, event.BuildDisplayCmd("error writing to save file")
			}
		}
	} else {
		imgString := process.RenderImageFile(m.controls.Settings, currentJob.sourcePath)
		if err := writeFile(currentJob.destinationPath, imgString); err != nil {
			return m, event.BuildDisplayCmd("error writing to save file")
		}
	}

	m.exportIndex += 1
	m.overwriteApproved = false // reset for next file
	displayMsg := fmt.Sprintf("%d/%d exports completed", m.exportIndex, len(m.exportQueue))
	displayCmd := event.BuildDisplayCmd(displayMsg)

	if m.exportIndex >= len(m.exportQueue) {
		m.waitingOnExport = false
		return m, displayCmd
	}

	return m, tea.Batch(event.StartRenderToExportCmd, displayCmd)
}

func writeFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	if _, err = w.WriteString(content); err != nil {
		return err
	}
	return w.Flush()
}

func (m Model) startExportingDir(msg event.StartExportMsg) (Model, tea.Cmd) {
	return m, event.BuildDisplayCmd(fmt.Sprintf("exporting %s", msg.SourcePath))
}

func (m Model) startExportingFile(msg event.StartExportMsg) (Model, tea.Cmd) {
	return m, event.BuildDisplayCmd(fmt.Sprintf("exporting %s", msg.SourcePath))
}

func (m Model) handleStartAdaptingMsg() (Model, tea.Cmd) {
	filename := m.controls.FileBrowser.SelectedFilename()
	message := fmt.Sprintf("generating palette from %s...", filename)
	return m, tea.Batch(event.BuildDisplayCmd(message), m.processAdaptingCmd)
}

func (m Model) handleFinishAdaptingMsg(msg event.FinishAdaptingMsg) (Model, tea.Cmd) {
	m.controls.Settings.Colors.PaletteControls.Adapter = m.controls.Settings.Colors.PaletteControls.Adapter.SetPalette(msg.Colors, msg.Name)
	return m, tea.Batch(event.StartRenderToViewCmd, event.BuildDisplayCmd("rendering..."))
}

type Foo struct {
	Bar string
}

func (m Model) handleLospecRequestMsg(msg event.LospecRequestMsg) (Model, tea.Cmd) {
	// make url request
	r, err := http.Get(msg.URL)
	if err != nil {
		return m, event.BuildDisplayCmd("error making lospec request")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return m, event.BuildDisplayCmd("error reading lospec response")
	}

	// parse json and populate LospecResponseMsg
	data := new(event.LospecData)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return m, event.BuildDisplayCmd("error decoding lospec request")
	}

	// build Data Cmd
	return m, event.BuildLospecResponseCmd(event.LospecResponseMsg{
		ID:   msg.ID,
		Page: msg.Page,
		Data: *data,
	})
}

func (m Model) handleLospecResponseMsg(msg event.LospecResponseMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.controls.Settings.Colors.PaletteControls.Lospec, cmd = m.controls.Settings.Colors.PaletteControls.Lospec.Update(msg)
	return m, cmd
}

func (m Model) processAdaptingCmd() tea.Msg {
	colors, name := adapt.GeneratePalette(m.controls.Settings.Colors.PaletteControls.Adapter, m.controls.FileBrowser.SelectedFile)
	return event.FinishAdaptingMsg{
		Name:   name,
		Colors: colors,
	}
}

func (m Model) handleDebug() (Model, tea.Cmd) {
	var b strings.Builder
	b.WriteString("=== DEBUG LOG ===\n\n")

	// App state
	b.WriteString(fmt.Sprintf("App: state=%d w=%d h=%d\n", m.state, m.w, m.h))
	b.WriteString(fmt.Sprintf("Active file: %s\n", m.controls.FileBrowser.ActiveFile))
	b.WriteString(fmt.Sprintf("Viewer: animating=%v waitingOnRender=%v\n\n", m.viewer.IsAnimating(), m.viewer.WaitingOnRender))

	// Controls state
	b.WriteString(m.controls.DebugState() + "\n")
	b.WriteString(m.controls.Settings.DebugState() + "\n\n")

	// Theme
	b.WriteString(fmt.Sprintf("Theme: %s (transparent=%v bg=%s)\n\n",
		style.ThemeNames[style.ActiveTheme.Name],
		style.ActiveTheme.Transparent,
		style.ColorHex(style.ActiveTheme.Bg)))

	// Viewer imgString
	b.WriteString("=== VIEWER IMG STRING ===\n")
	b.WriteString(m.viewer.ImgString())
	b.WriteString("\n=== END VIEWER IMG STRING ===\n\n")

	// Full ANSI output
	b.WriteString("=== FULL ANSI OUTPUT ===\n")
	b.WriteString(m.View().Content)
	b.WriteString("\n=== END ===\n")

	_ = os.WriteFile("debug.log", []byte(b.String()), 0644)
	return m, event.BuildDisplayCmd("debug.log written")
}

func (m Model) handleControlsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.controls, cmd = m.controls.Update(msg)

	// Check if the controls wants to quit
	if m.controls.ShouldQuit {
		m.controls.ShouldQuit = false
		m.quitConfirm = true
		m.quitFocusYes = true
		return m, nil
	}

	// Check if the controls wants us to open a modal
	if m.controls.OpenModal != nil {
		result := m.controls.OpenModal
		m.controls.OpenModal = nil
		switch result.Action {
		case filemenu.LoadSettings:
			m = m.openLoadModal()
		case filemenu.SaveSettings:
			m = m.openSaveModal()
		case filemenu.Export:
			switch result.ExportSub {
			case filemenu.ExportCurrentFile:
				m = m.openExportFileModal()
			case filemenu.ExportBatchProcess:
				m = m.openExportBatchModal()
			}
		case filemenu.Preferences:
			m = m.openPrefsModal()
		}
	}

	m.savePrefs()
	return m, cmd
}

func (m *Model) savePrefs() {
	dirs := prefs.Dirs{
		Browse:       m.controls.FileBrowser.SelectedDir,
		ExportSource: m.controls.Export.Source.Browser.SelectedDir,
		ExportDest:   m.controls.Export.Destination.GetSelected(),
		SaveDir:      m.controls.Settings.SaveLoad.SelectedDirPath(),
		LoadFile:     m.controls.Settings.SaveLoad.SelectedLoadDir(),
		PaletteLoad:  m.controls.Settings.Colors.PaletteControls.Loader.FileBrowser.SelectedDir,
	}
	if dirs != m.prefs.Dirs {
		m.prefs.Dirs = dirs
		m.prefs.Save()
	}
}

func (m Model) saveLatestSettings() {
	cfg := m.controls.Settings.ExportConfig()
	_ = settings.SaveConfig(cfg, prefs.LatestSettingsPath())
}

func (m Model) handleDisplayMsg(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.display, cmd = m.display.Update(msg)
	return m, cmd
}

func (m Model) handleCopy() (Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.viewer.View()); err != nil {
		return m, event.BuildDisplayCmd("Error copying to clipboard")
		// we should have a place in the UI where we display errors or processing messages,
		// and package our desired event to the user in a specific command
	}
	filename := m.controls.FileBrowser.SelectedFilename()
	name := strings.Split(filename, ".")[0] // strip extension
	return m, event.BuildDisplayCmd(fmt.Sprintf("copied %s to clipboard", name))
}

func (m Model) handleSave() (Model, tea.Cmd) {
	name := strings.Split(m.controls.FileBrowser.SelectedFilename(), ".")[0]

	if m.viewer.IsAnimating() {
		frames := m.viewer.Frames()
		if err := os.MkdirAll(name, 0755); err != nil {
			return m, event.BuildDisplayCmd("error creating save directory")
		}
		for i, frame := range frames {
			filename := fmt.Sprintf("%s_%03d.ansi", name, i+1)
			filePath := fmt.Sprintf("%s/%s", name, filename)
			file, err := os.Create(filePath)
			if err != nil {
				return m, event.BuildDisplayCmd(fmt.Sprintf("error creating frame file %s", filename))
			}
			w := bufio.NewWriter(file)
			_, err = w.WriteString(frame)
			if err != nil {
				file.Close()
				return m, event.BuildDisplayCmd(fmt.Sprintf("error writing frame %s", filename))
			}
			w.Flush()
			file.Close()
		}
		return m, event.BuildDisplayCmd(fmt.Sprintf("saved %d frames to %s/", len(frames), name))
	}

	filename := fmt.Sprintf("%s.ansi", name)
	file, err := os.Create(filename)
	if err != nil {
		return m, event.BuildDisplayCmd("error creating save file")
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString(m.viewer.View())
	if err != nil {
		return m, event.BuildDisplayCmd("error writing to save file")
	}
	w.Flush()

	return m, event.BuildDisplayCmd(fmt.Sprintf("saved to %s", filename))
}
