package state

import (
	"fmt"
	"os"
)

// Model contains all info used to render the full app
type Model struct {
	Browser *Browser

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
	// target Directory
	// included filetypes
	// export Directory
	// filename extension
	// save video frames to subdirectory (bool)
}

type Browser struct {
	Directory string
	Filename  string
}

func New() *Model {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current Directory:", err)
		os.Exit(1)
	}

	return &Model{
		Browser: &Browser{
			Directory: dir,
			Filename:  "",
		},
	}
}
