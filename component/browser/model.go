package browser

import (
	"fmt"
	"github.com/Zebbeni/ansizalizer/component/item"
	"github.com/Zebbeni/ansizalizer/component/keyboard"
	"github.com/Zebbeni/ansizalizer/component/menu"
	"github.com/Zebbeni/ansizalizer/state"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
)

var (
	acceptedFileExts = []string{".jpg", ".mod"}
)

type Model struct {
	state  *state.Browser
	menu   *menu.Model
	keymap *keyboard.Map
}

// build a list of all directories and image files at current location
// current directory
// subdirectories (alphabetical)
// image files (alphabetical)

func New(state *state.Browser, k *keyboard.Map) *Model {
	m := menu.New(nil, k)

	model := &Model{
		state: state,
		menu:  m,
	}

	model.updateItems()

	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return m.menu.View()
}

func (m *Model) HandleKeyMsg(msg tea.KeyMsg) bool {
	return m.menu.HandleKeyMsg(msg)
}

func (m *Model) updateDirectory(b *state.Browser, dir string) {
	b.Directory = dir
	m.updateItems()
}

func (m *Model) updateFile(b *state.Browser, file string) {
	b.Filename = file
}

func (m *Model) updateItems() {
	entries, err := os.ReadDir(m.state.Directory)
	if err != nil {
		fmt.Println("Unable to read directory entries", err)
		os.Exit(1)
	}

	dirItems := make([]item.Model, 1)
	dirName := fmt.Sprintf("../%s", filepath.Base(m.state.Directory))
	dirParentName := filepath.Dir(m.state.Directory)
	dirItems[0] = item.New(dirName, func() {
		m.updateDirectory(m.state, dirParentName)
	})

	fileItems := make([]item.Model, 0)

	for _, entry := range entries {
		fullPath := filepath.Join(m.state.Directory, entry.Name())
		if entry.IsDir() {
			dirItems = append(dirItems, item.New(entry.Name(), func() {
				m.updateDirectory(m.state, fullPath)
			}))
			continue
		}
		for _, ext := range acceptedFileExts {
			if filepath.Ext(entry.Name()) == ext {
				fileItems = append(fileItems, item.New(entry.Name(), func() {
					m.updateFile(m.state, fullPath)
				}))
				continue
			}
		}
	}

	m.menu.Items = append(dirItems, fileItems...)
}
