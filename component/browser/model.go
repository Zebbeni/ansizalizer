package browser

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/menu"
)

var (
	acceptedFileExts = []string{".jpg", ".png"}
)

type Model struct {
	navState State
	menu     menu.Model

	SetDir  SetDirPath
	SetFile SetFilePath
}

func New(setDir SetDirPath, setFile SetFilePath) Model {
	m := menu.New(nil)

	model := Model{
		menu:    m,
		SetDir:  setDir,
		SetFile: setFile,
	}

	model.updateItems()

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.SetDir(m.navState.DirPath)
	m.SetFile(m.navState.Filepath)

	return m.menu.Update(msg)
}

func (m Model) View() string {
	return m.menu.View()
}

func (m Model) GetActivePosition() (float64, float64) {
	return m.menu.GetActivePosition()
}

func (m Model) updateDirectory(dir string) {
	m.navState.DirPath = dir
	m.SetDir(dir)
	m.updateItems()
}

func (m Model) updateImageFile(file string) {
	if m.navState.Filepath == file {
		return
	}

	m.navState.Filepath = file
}

func (m Model) updateItems() {
	dirPath := m.navState.DirPath

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Unable to read directory entries", err)
		os.Exit(1)
	}

	// initialize directory items with an item that allows going up a directory
	dirBaseName := fmt.Sprintf("../%s", filepath.Base(dirPath))
	dirParentName := filepath.Dir(dirPath)
	dirItems := []item.Model{
		item.New(dirBaseName, func() {
			m.updateDirectory(dirParentName)
		}),
	}
	fileItems := make([]item.Model, 0)

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			dirItems = append(dirItems, item.New(entry.Name(), func() {
				m.updateDirectory(fullPath)
			}))
			continue
		}
		for _, ext := range acceptedFileExts {
			if filepath.Ext(entry.Name()) == ext {
				fileItems = append(fileItems, item.New(entry.Name(), func() {
					m.updateImageFile(fullPath)
				}))
				continue
			}
		}
	}

	m.menu.Items = append(dirItems, fileItems...)
	m.menu.ResetActive()
}
