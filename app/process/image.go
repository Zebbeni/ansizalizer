package process

import (
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
