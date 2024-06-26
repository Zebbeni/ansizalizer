package process

import (
	"bufio"
	"image"
	"os"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/controls/settings"
	"github.com/Zebbeni/ansizalizer/controls/settings/characters"
)

var (
	black = colorful.Color{}
)

func RenderImageFile(s settings.Model, imgFilePath string) string {
	if imgFilePath == "" {
		return "Browse an image to render"
	}

	var img image.Image
	imgFile, err := os.Open(imgFilePath)
	if err != nil {
		return "Could not open image " + imgFilePath
	}
	defer imgFile.Close()
	imageReader := bufio.NewReader(imgFile)
	img, _, err = image.Decode(imageReader)
	if err != nil {
		return "Could not decode image " + imgFilePath
	}

	renderer := New(s)
	imgString := renderer.process(img)
	return imgString
}

func (m Renderer) process(input image.Image) string {
	isTrueColor, _, palette := m.Settings.Colors.GetSelected()
	if !isTrueColor && len(palette.Colors()) == 0 {
		return "Choose a color palette"
	}
	mode, _, _, _ := m.Settings.Characters.Selected()
	switch mode {
	case characters.Ascii:
		return m.processAscii(input)
	case characters.Unicode:
		return m.processUnicode(input)
	case characters.Custom:
		return m.processCustom(input)
	}
	return "Choose a character type"
}
