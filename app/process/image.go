package process

import (
	"strings"
	"time"

	"github.com/Zebbeni/ansiart"
	"github.com/Zebbeni/ansizalizer/controls/settings"
)

func RenderImageFile(s settings.Model, imgFilePath string) string {
	if imgFilePath == "" {
		return "Browse an image to render"
	}

	opts := settingsToOptions(s)
	result, err := ansiart.RenderFile(imgFilePath, opts)
	if err != nil {
		return err.Error()
	}
	return result
}

func RenderGIFFile(s settings.Model, imgFilePath string) ([]string, []time.Duration) {
	opts := settingsToOptions(s)
	frames, err := ansiart.RenderGIF(imgFilePath, opts)
	if err != nil {
		return []string{err.Error()}, []time.Duration{0}
	}
	rendered := make([]string, len(frames))
	delays := make([]time.Duration, len(frames))
	for i, f := range frames {
		rendered[i] = f.Content
		delays[i] = f.Delay
	}
	return rendered, delays
}

func IsGIFFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".gif")
}
