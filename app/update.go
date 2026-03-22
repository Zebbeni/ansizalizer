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
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/app/adapt"
	"github.com/Zebbeni/ansizalizer/app/process"
	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleStartRenderToViewCmd() (Model, tea.Cmd) {
	m.viewer.WaitingOnRender = true
	return m, m.processRenderToViewCmd
}

func (m Model) handleFinishRenderToViewMsg(msg event.FinishRenderToViewMsg) (Model, tea.Cmd) {
	// cut out early if the finished render is for a previously selected image
	if msg.FilePath != m.controls.FileBrowser.ActiveFile {
		return m, nil
	}
	var cmd tea.Cmd
	m.controls.Settings.Alpha.AlphaImage = false
	re := regexp.MustCompile(`(?i)\.(png|gif)$`)
	if re.Match([]byte(m.controls.FileBrowser.ActiveFile)) {
		m.controls.Settings.Alpha.AlphaImage = true
	}
	m.viewer, cmd = m.viewer.Update(msg)
	return m, cmd
}

func (m Model) processRenderToViewCmd() tea.Msg {
	re := regexp.MustCompile(`(?i)\.(png|gif)$`)
	filePath := m.controls.FileBrowser.ActiveFile

	colorsString := "true color"
	alphaString := "no alpha channel"
	if m.controls.Settings.Colors.IsLimited() {
		palette := m.controls.Settings.Colors.GetCurrentPalette()
		colorsString = palette.Title()
	}
	if m.controls.Settings.Alpha.ShouldOutputAlpha() && re.Match([]byte(filePath)) {
		alphaString = "alpha channel"
	}

	// Animated GIF path
	if process.IsGIFFile(filePath) {
		frames, _ := process.RenderGIFFile(m.controls.Settings, filePath)
		if len(frames) > 1 {
			return event.FinishRenderGIFToViewMsg{
				FilePath:     filePath,
				Frames:       frames,
				Delay:        m.controls.Settings.Animation.Delay(),
				ColorsString: colorsString,
				AlphaString:  alphaString,
			}
		}
		// Single-frame GIF: use normal path
		return event.FinishRenderToViewMsg{
			FilePath:     filePath,
			ImgString:    frames[0],
			ColorsString: colorsString,
			AlphaString:  alphaString,
		}
	}

	// Non-GIF: existing path
	imgString := process.RenderImageFile(m.controls.Settings, filePath)
	return event.FinishRenderToViewMsg{
		FilePath:     filePath,
		ImgString:    imgString,
		ColorsString: colorsString,
		AlphaString:  alphaString,
	}
}

func (m Model) handleFinishRenderGIFToViewMsg(msg event.FinishRenderGIFToViewMsg) (Model, tea.Cmd) {
	if msg.FilePath != m.controls.FileBrowser.ActiveFile {
		return m, nil
	}
	m.controls.Settings.Alpha.AlphaImage = true
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

	return m, tea.Batch(event.StartRenderToExportCmd, event.BuildDisplayCmd(fmt.Sprintf("%d jobs queued", len(exportQueue))))
}

func (m Model) handleRenderToExportMsg() (Model, tea.Cmd) {

	currentJob := m.exportQueue[m.exportIndex]

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
	filename := m.controls.FileBrowser.ActiveFilename()
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
	colors, name := adapt.GeneratePalette(m.controls.Settings.Colors.PaletteControls.Adapter, m.controls.FileBrowser.ActiveFile)
	return event.FinishAdaptingMsg{
		Name:   name,
		Colors: colors,
	}
}

func (m Model) handleControlsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.controls, cmd = m.controls.Update(msg)
	return m, cmd
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
	filename := m.controls.FileBrowser.ActiveFilename()
	name := strings.Split(filename, ".")[0] // strip extension
	return m, event.BuildDisplayCmd(fmt.Sprintf("copied %s to clipboard", name))
}

func (m Model) handleSave() (Model, tea.Cmd) {
	name := strings.Split(m.controls.FileBrowser.ActiveFilename(), ".")[0]

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
