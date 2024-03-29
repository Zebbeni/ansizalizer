package browser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

type item struct {
	name  string
	path  string
	isDir bool
	isTop bool
}

func (i item) FilterValue() string {
	return i.name
}

func (i item) Title() string {
	if i.isTop {
		return "↑"
	}
	if i.isDir {
		return fmt.Sprintf("%s/", i.name)
	}
	return i.name
}

func (i item) Description() string {
	if i.isDir {
		return "directory"
	}
	return "file"
}

func getItems(extensions map[string]bool, dir string) []list.Item {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory entries:", err)
		os.Exit(1)
	}

	parentPath := filepath.Dir(dir)
	parentName := filepath.Base(parentPath)
	parentItem := item{name: parentName, path: parentPath, isDir: true, isTop: true}

	dirItems := []list.Item{parentItem}
	fileItems := make([]list.Item, 0)

	for _, e := range entries {
		path := fmt.Sprintf("%s/%s", dir, e.Name())

		if e.IsDir() {
			name := e.Name()
			dirItem := item{name: name, path: path, isDir: true, isTop: false}
			dirItems = append(dirItems, dirItem)
			continue
		}

		ext := filepath.Ext(e.Name())
		if _, ok := extensions[ext]; ok {
			fileItem := item{name: e.Name(), path: path, isDir: false, isTop: false}
			fileItems = append(fileItems, fileItem)
		}
	}

	return append(dirItems, fileItems...)
}
