package process

import (
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansipx"
	"github.com/Zebbeni/ansizalizer/controls/settings"
)

func RenderImageFile(s settings.Model, imgFilePath string) string {
	return RenderImageFileWithBg(s, imgFilePath, nil)
}

func RenderImageFileWithBg(s settings.Model, imgFilePath string, solidBg *colorful.Color) string {
	if imgFilePath == "" {
		return "Browse an image to render"
	}

	opts := settingsToOptionsWithBg(s, solidBg)
	result, err := ansipx.RenderFile(imgFilePath, opts)
	if err != nil {
		return err.Error()
	}
	return result
}

func RenderGIFFile(s settings.Model, imgFilePath string) ([]string, []time.Duration) {
	return RenderGIFFileWithBg(s, imgFilePath, nil)
}

func RenderGIFFileWithBg(s settings.Model, imgFilePath string, solidBg *colorful.Color) ([]string, []time.Duration) {
	opts := settingsToOptionsWithBg(s, solidBg)
	frames, err := ansipx.RenderGIF(imgFilePath, opts)
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
