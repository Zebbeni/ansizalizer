package modal

import (
	"fmt"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
	"github.com/Zebbeni/ansizalizer/global"
	"github.com/Zebbeni/ansizalizer/style"
)

const (
	Width       = 50
	Height      = 20
	LoadWidth   = 60
	LoadHeight  = 24
	BatchWidth  = 106
	BatchHeight = 30
)

type Kind int

const (
	LoadKind Kind = iota
	SaveKind
	ExportFileKind
	ExportBatchKind
	PrefsKind
)

type Focus int

const (
	FocusBrowser Focus = iota
	FocusLoadPanel  // load browser panel focused but not active
	FocusSavePanel   // save browser panel focused but not active
	FocusSaveSelect  // "Select <dir>" button above save browser
	FocusSaveBrowser // browsing inside save panel
	FocusSaveInput   // filename input in save modal
	FocusInput
	FocusConfirm
	FocusCancel
	FocusSourcePanel   // source panel focused (not browsing)
	FocusSourceSelect  // "Select <dir>?" button above source browser
	FocusSourceBrowser // browsing inside source panel
	FocusSubDirsYes
	FocusSubDirsNo
	FocusDestPanel     // dest panel focused (not browsing)
	FocusDestSelect    // "Select <dir>?" button above dest browser
	FocusDestBrowser   // browsing inside dest panel
	FocusPrefsShowHelp
	FocusPrefsRestoreTheme
	FocusPrefsRestoreSettings
	FocusPrefsSkipSplash
	FocusPrefsClose
)

// Result is set when the user confirms an action.
type Result struct {
	Kind            Kind
	Path            string // full file path for load or save
	SourcePath      string // for export
	DestinationPath string // for export
	IsDir           bool   // for export batch
	UseSubDirs      bool   // for export batch
	ShowHelp        bool   // for preferences
	RestoreTheme    bool   // for preferences
	RestoreSettings bool   // for preferences
	SplashOnStart   bool   // for preferences
}

type Model struct {
	Kind    Kind
	Open    bool
	focus   Focus
	browser browser.Model
	input   textinput.Model
	result  *Result

	dir   string // current directory for save
	width int

	// Export-specific fields
	sourceFile  string       // source file path for ExportFileKind
	sourceBrowser browser.Model // source directory browser for ExportBatchKind
	destBrowser   browser.Model // destination directory browser
	useSubDirs    bool          // include subdirectories in batch export
	confirmedSrc  string        // confirmed source directory
	confirmedDest string        // confirmed destination directory
	imageCount    int           // number of processable images in confirmed source

	// Preferences fields
	prefsShowHelp        bool
	prefsRestoreTheme    bool
	prefsRestoreSettings bool
	prefsSplashOnStart   bool
}

var jsonExtensions = map[string]bool{".json": true}

func NewLoad(settingsDir string) Model {
	contentArea := LoadWidth - 6 // 2 border + 4 padding
	panelContent := contentArea - 2 // panel border
	browserW := panelContent - 6 // extra margin for emoji width
	b := browser.NewAtDir(settingsDir, jsonExtensions, browserW)
	b.SetActive(false)
	return Model{
		Kind:    LoadKind,
		Open:    true,
		focus:   FocusLoadPanel,
		browser: b,
		width:   LoadWidth,
	}
}

func NewSave(settingsDir string) Model {
	contentArea := LoadWidth - 6
	panelContent := contentArea - 2
	browserW := panelContent - 6 // extra margin for emoji width

	b := browser.NewAtDir(settingsDir, nil, browserW)
	b.SetActive(false)

	input := textinput.New()
	input.Placeholder = "settings"
	input.Prompt = "Filename: "
	input.CharLimit = 30
	input.SetWidth(contentArea - 14)
	return Model{
		Kind:    SaveKind,
		Open:    true,
		focus:   FocusSavePanel,
		input:   input,
		browser: b,
		dir:     settingsDir,
		width:   LoadWidth,
	}
}

func NewExportFile(sourceFile, destDir string) Model {
	contentArea := LoadWidth - 6
	panelContent := contentArea - 2
	browserW := panelContent - 6
	destBrowser := browser.NewAtDir(destDir, nil, browserW)
	destBrowser.SetActive(false)
	return Model{
		Kind:        ExportFileKind,
		Open:        true,
		focus:       FocusDestPanel,
		destBrowser: destBrowser,
		sourceFile:  sourceFile,
		confirmedDest: destDir,
		width:       LoadWidth,
	}
}

func NewExportBatch(sourceDir, destDir string) Model {
	panelW := (BatchWidth - 9) / 2 // each panel's content width
	sourceBrowser := browser.NewAtDir(sourceDir, nil, panelW)
	sourceBrowser.SetActive(false)
	destBrowser := browser.NewAtDir(destDir, nil, panelW)
	destBrowser.SetActive(false)
	return Model{
		Kind:          ExportBatchKind,
		Open:          true,
		focus:         FocusSourcePanel,
		sourceBrowser: sourceBrowser,
		destBrowser:   destBrowser,
		width:         BatchWidth,
	}
}

func NewPrefs(showHelp, restoreTheme, restoreSettings, splashOnStart bool) Model {
	return Model{
		Kind:                 PrefsKind,
		Open:                 true,
		focus:                FocusPrefsShowHelp,
		width:                Width,
		prefsShowHelp:        showHelp,
		prefsRestoreTheme:    restoreTheme,
		prefsRestoreSettings: restoreSettings,
		prefsSplashOnStart:   splashOnStart,
	}
}

func (m Model) GetResult() *Result {
	return m.result
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.Open {
		return m, nil
	}

	switch m.Kind {
	case LoadKind:
		return m.updateLoad(msg)
	case SaveKind:
		return m.updateSave(msg)
	case PrefsKind:
		return m.updatePrefs(msg)
	case ExportFileKind:
		return m.updateExportFile(msg)
	case ExportBatchKind:
		return m.updateExportBatch(msg)
	}
	return m, nil
}

func (m Model) updateLoad(msg tea.Msg) (Model, tea.Cmd) {
	if m.focus == FocusBrowser {
		prevFile := m.browser.SelectedFile
		var cmd tea.Cmd
		m.browser, cmd = m.browser.Update(msg)

		if m.browser.ShouldClose {
			m.browser.ShouldClose = false
			m.browser.SetActive(false)
			m.focus = FocusLoadPanel
		} else if m.browser.SelectedFile != prevFile && m.browser.SelectedFile != "" {
			// File selected — move to confirm
			m.browser.SetActive(false)
			m.focus = FocusConfirm
		}
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		switch m.focus {
		case FocusLoadPanel:
			m.focus = FocusBrowser
			m.browser.SetActive(true)
		case FocusConfirm:
			if m.browser.SelectedFile != "" {
				m.result = &Result{Kind: LoadKind, Path: m.browser.SelectedFile}
				m.Open = false
			}
		case FocusCancel:
			m.Open = false
		}
	case key.Matches(keyMsg, event.KeyMap.Left):
		if m.focus == FocusCancel {
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Right):
		if m.focus == FocusConfirm {
			m.focus = FocusCancel
		}
	case key.Matches(keyMsg, event.KeyMap.Down):
		switch m.focus {
		case FocusLoadPanel:
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusConfirm, FocusCancel:
			m.focus = FocusLoadPanel
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		m.Open = false
	}
	return m, nil
}

func (m Model) updateSave(msg tea.Msg) (Model, tea.Cmd) {
	// Browser active
	if m.focus == FocusSaveBrowser {
		var cmd tea.Cmd
		m.browser, cmd = m.browser.Update(msg)
		if m.browser.ShouldClose {
			m.browser.ShouldClose = false
			m.browser.SetActive(false)
			m.focus = FocusSaveSelect
		}
		// Track current directory for the Select button
		m.dir = m.browser.SelectedDir
		return m, cmd
	}

	// Filename input active
	if m.focus == FocusSaveInput && m.input.Focused() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, event.KeyMap.Enter):
				m.input.Blur()
				m.focus = FocusConfirm
				return m, nil
			case key.Matches(keyMsg, event.KeyMap.Esc):
				m.input.Blur()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case key.Matches(keyMsg, event.KeyMap.Enter):
		switch m.focus {
		case FocusSavePanel:
			// Expand the panel — focus on Select button
			m.focus = FocusSaveSelect
		case FocusSaveSelect:
			// Confirm directory, collapse panel, move to filename
			m.dir = m.browser.SelectedDir
			m.focus = FocusSaveInput
		case FocusSaveInput:
			m.input.Focus()
		case FocusConfirm:
			filename := m.input.Value()
			if filename == "" {
				filename = "settings"
			}
			m.result = &Result{Kind: SaveKind, Path: filepath.Join(m.dir, filename+".json")}
			m.Open = false
		case FocusCancel:
			m.Open = false
		}
	case key.Matches(keyMsg, event.KeyMap.Left):
		if m.focus == FocusCancel {
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Right):
		if m.focus == FocusConfirm {
			m.focus = FocusCancel
		}
	case key.Matches(keyMsg, event.KeyMap.Down):
		switch m.focus {
		case FocusSavePanel:
			m.focus = FocusSaveInput
		case FocusSaveSelect:
			m.focus = FocusSaveBrowser
			m.browser.SetActive(true)
		case FocusSaveInput:
			m.focus = FocusConfirm
		}
	case key.Matches(keyMsg, event.KeyMap.Up):
		switch m.focus {
		case FocusSaveSelect:
			// Collapse back to panel
			m.focus = FocusSavePanel
		case FocusSaveInput:
			m.focus = FocusSavePanel
		case FocusConfirm, FocusCancel:
			m.focus = FocusSaveInput
		}
	case key.Matches(keyMsg, event.KeyMap.Esc):
		switch m.focus {
		case FocusSaveSelect:
			// Collapse panel
			m.focus = FocusSavePanel
		default:
			m.Open = false
		}
	}
	return m, nil
}

func (m Model) View() string {
	w, h := Width, Height
	if m.Kind == LoadKind || m.Kind == SaveKind || m.Kind == ExportFileKind {
		w, h = LoadWidth, LoadHeight
	} else if m.Kind == ExportBatchKind {
		w, h = BatchWidth, BatchHeight
	}

	var title, body string

	switch m.Kind {
	case LoadKind:
		title = "Load Settings"
		body = m.loadView()
	case SaveKind:
		title = "Save Settings"
		body = m.saveView()
	case PrefsKind:
		title = "Preferences"
		body = m.prefsView()
	case ExportFileKind:
		title = "Export Current File"
		body = m.exportFileView()
	case ExportBatchKind:
		title = "Batch Export"
		body = m.exportBatchView()
	}

	buttons := ""
	if m.Kind != PrefsKind {
		buttons = m.renderButtons()
	}

	if m.Kind == ExportBatchKind {
		// Use BoxWithLabel for the outer modal
		content := lipgloss.JoinVertical(lipgloss.Center, body, "", buttons)

		textStyle := style.BgStyle().
			AlignHorizontal(lipgloss.Center).
			Padding(0, 1).
			Foreground(style.SelectedColor1)
		borderStyle := style.BgStyle().
			Border(style.HeavyBorder()).
			BorderForeground(style.SelectedColor1).
			BorderBackground(style.ActiveTheme.Bg)

		renderer := style.BoxWithLabel{
			BoxStyle:   borderStyle,
			LabelStyle: textStyle,
		}
		return style.ApplyBg(renderer.Render(title, content, w-2), 0)
	}

	// Standard modal for load/save/export-file/prefs
	titleRendered := style.SelectedTitle.
		Width(w - 6).
		AlignHorizontal(lipgloss.Center).
		Render(title)

	content := lipgloss.JoinVertical(lipgloss.Center, titleRendered, "", body, "", buttons)

	border := style.BgStyle().
		Width(w).
		Height(h).
		Border(style.HeavyBorder()).
		BorderForeground(style.SelectedColor1).
		BorderBackground(style.ActiveTheme.Bg).
		Padding(1, 2).
		AlignHorizontal(lipgloss.Center)

	if m.Kind == PrefsKind {
		border = border.AlignVertical(lipgloss.Center)
	}

	return style.ApplyBg(border.Render(content), 0)
}

func (m Model) loadView() string {
	contentArea := LoadWidth - 6 // outer border + padding
	panelContent := contentArea - 2 // inner panel border

	desc := style.DimmedTitle.Italic(true).
		Width(contentArea).
		AlignHorizontal(lipgloss.Center).
		Render("Select a .json settings file to load.")

	panelFocused := m.focus == FocusLoadPanel
	panelActive := m.focus == FocusBrowser

	renderColor := style.DimmedColor1
	if panelActive {
		renderColor = style.NormalColor1
	} else if panelFocused {
		renderColor = style.SelectedColor1
	}

	border := lipgloss.RoundedBorder()
	if panelFocused {
		border = style.HeavyBorder()
	}
	borderStyle := style.BgStyle().
		Border(border).
		BorderForeground(renderColor).
		BorderBackground(style.ActiveTheme.Bg)
	textSt := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1).
		Foreground(renderColor)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: textSt,
	}

	content := lipgloss.NewStyle().PaddingTop(1).Render(m.browser.View())
	title := "Browse " + filepath.Base(m.browser.SelectedDir) + "/"
	panel := renderer.Render(title, content, panelContent)

	// Show selected file below the panel
	selected := ""
	if m.browser.SelectedFile != "" {
		selected = style.NormalTitle.
			Width(contentArea).
			AlignHorizontal(lipgloss.Center).
			Render(filepath.Base(m.browser.SelectedFile))
	}

	return lipgloss.JoinVertical(lipgloss.Center, desc, "", panel, selected)
}

func (m Model) saveView() string {
	contentArea := LoadWidth - 6
	panelContent := contentArea - 2

	desc := style.DimmedTitle.Italic(true).
		Width(contentArea).
		AlignHorizontal(lipgloss.Center).
		Render("Choose a directory and filename to save settings.")

	// Destination panel — collapsed or expanded
	panelExpanded := m.focus == FocusSaveSelect || m.focus == FocusSaveBrowser
	panelFocused := m.focus == FocusSavePanel
	panelActive := m.focus == FocusSaveBrowser || m.focus == FocusSaveSelect

	renderColor := style.DimmedColor1
	if panelActive {
		renderColor = style.NormalColor1
	} else if panelFocused {
		renderColor = style.SelectedColor1
	}

	border := lipgloss.RoundedBorder()
	if panelFocused {
		border = style.HeavyBorder()
	}
	borderStyle := style.BgStyle().
		Border(border).
		BorderForeground(renderColor).
		BorderBackground(style.ActiveTheme.Bg)
	titleSt := style.BgStyle().
		AlignHorizontal(lipgloss.Center).
		Padding(0, 1).
		Foreground(renderColor)

	renderer := style.BoxWithLabel{
		BoxStyle:   borderStyle,
		LabelStyle: titleSt,
	}

	var panel string
	if panelExpanded {
		selectBtn := m.renderSelectButton(m.browser.SelectedDir, m.focus == FocusSaveSelect)
		browseContent := lipgloss.JoinVertical(lipgloss.Left, selectBtn, m.browser.View())
		panel = renderer.Render("Destination", browseContent, panelContent)
	} else {
		summary := style.DimmedTitle.Italic(true).
			Width(panelContent).
			AlignHorizontal(lipgloss.Center).
			Padding(0, 1).
			Render(filepath.Base(m.dir) + "/")
		panel = renderer.Render("Destination", summary, panelContent)
	}

	// Filename input
	inputFocused := m.focus == FocusSaveInput
	bg := style.BgStyle()
	promptSt := bg.Foreground(style.DimmedColor1)
	textSt := bg.Foreground(style.DimmedColor1)
	cursorColor := style.DimmedColor1
	if m.input.Focused() {
		promptSt = bg.Foreground(style.SelectedColor1)
		textSt = bg.Foreground(style.SelectedColor1)
		cursorColor = style.SelectedColor1
	} else if inputFocused {
		promptSt = bg.Foreground(style.NormalColor1)
		textSt = bg.Foreground(style.NormalColor1)
		cursorColor = style.NormalColor1
	}
	styles := m.input.Styles()
	placeholderSt := bg.Foreground(style.ExtraDimColor)
	styles.Focused.Prompt = promptSt.Width(12)
	styles.Focused.Text = textSt
	styles.Focused.Placeholder = placeholderSt
	styles.Blurred.Prompt = promptSt.Width(12)
	styles.Blurred.Text = textSt
	styles.Blurred.Placeholder = placeholderSt
	styles.Cursor.Blink = true
	styles.Cursor.Color = cursorColor
	m.input.SetStyles(styles)
	m.input.SetVirtualCursor(m.input.Focused())
	filenameRow := m.input.View()

	return lipgloss.JoinVertical(lipgloss.Center, desc, "", panel, "", filenameRow)
}

func (m Model) renderButtons() string {
	confirmLabel := "  Load  "
	if m.Kind == LoadKind && m.browser.SelectedFile != "" {
		confirmLabel = " Load " + filepath.Base(m.browser.SelectedFile) + " "
	}
	switch m.Kind {
	case SaveKind:
		filename := m.input.Value()
		if filename == "" {
			filename = "settings"
		}
		confirmLabel = " Save " + filepath.Base(m.dir) + "/" + filename + ".json "
	case ExportFileKind:
		confirmLabel = " Export " + filepath.Base(m.confirmedDest) + "/" + m.exportFileName() + " "
	case ExportBatchKind:
		confirmLabel = fmt.Sprintf(" Process (%d Images) ", m.imageCount)
	}

	confirmStyle := style.NormalButton
	if m.focus == FocusConfirm {
		confirmStyle = style.FocusButton
	}
	cancelStyle := style.NormalButton
	if m.focus == FocusCancel {
		cancelStyle = style.FocusButton
	}

	if m.Kind == ExportFileKind {
		return lipgloss.JoinVertical(lipgloss.Center,
			confirmStyle.Render(confirmLabel),
			"",
			cancelStyle.Render(" Cancel "),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center,
		confirmStyle.Render(confirmLabel),
		"  ",
		cancelStyle.Render(" Cancel "),
	)
}

// countImages counts image files in a directory, optionally including subdirectories.
func countImages(dir string, useSubDirs bool) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			if useSubDirs {
				count += countImages(filepath.Join(dir, e.Name()), true)
			}
			continue
		}
		ext := filepath.Ext(e.Name())
		if _, ok := global.ImgExtensions[ext]; ok {
			count++
		}
	}
	return count
}
