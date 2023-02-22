package state

import (
	"fmt"
	"os"
)

// Model contains all info used to render the full app
type Model struct {
	// 'File Browser' state
	// current directory
	directory string
	// current filename (changing this should trigger a new preview)

	// 'Settings' state
	// sampling method (changing any of these should trigger a new preview)
	// character set
	// palette
	// max width
	// max height

	// 'Viewer' state
	// []image (len > 1 means we're viewing a video)
	// frame duration

	// 'Process' state
	// target file
	// target directory
	// included filetypes
	// export directory
	// filename extension
	// save video frames to subdirectory (bool)
}

func New() *Model {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		os.Exit(1)
	}

	return &Model{
		directory: dir,
	}
}
