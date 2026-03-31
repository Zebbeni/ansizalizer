package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/list"

	"github.com/Zebbeni/ansizalizer/debug"
	"github.com/Zebbeni/ansizalizer/style"
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
		return "⮤ 🗀  " + style.BgStyle().Italic(true).Render(fmt.Sprintf("%s/", i.name))
	}
	if i.isDir {
		return fmt.Sprintf("🗀  %s/", i.name)
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
		path := filepath.Join(dir, e.Name())

		if e.IsDir() {
			name := e.Name()
			dirItem := item{name: name, path: path, isDir: true, isTop: false}
			dirItems = append(dirItems, dirItem)
			continue
		}

		ext := filepath.Ext(e.Name())
		if _, ok := extensions[strings.ToLower(ext)]; ok {
			fileItem := item{name: e.Name(), path: path, isDir: false, isTop: false}
			fileItems = append(fileItems, fileItem)
		}
	}

	debug.Log("getItems(%s): %d dirs, %d files, %d total entries from os.ReadDir",
		dir, len(dirItems)-1, len(fileItems), len(entries))
	for _, e := range entries {
		debug.Log("  entry: %s (isDir=%v)", e.Name(), e.IsDir())
	}

	return append(dirItems, fileItems...)
}
