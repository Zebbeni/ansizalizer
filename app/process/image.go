package process

import (
	"bufio"
	"image"
	"os"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/Zebbeni/ansizalizer/controls/options"
	"github.com/Zebbeni/ansizalizer/controls/options/characters"
)

const PROPORTION = 0.44

var (
	black = colorful.Color{}
)

func RenderImageFile(s options.Model, imgFilePath string) string {
	if imgFilePath == "" {
		return "Choose an image to render"
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
	if m.Settings.Colors.IsLimited() && len(m.Settings.Colors.GetCurrentPalette()) == 0 {
		return "Choose a description"
	}
	mode, _, _ := m.Settings.Characters.Selected()
	if mode == characters.Ascii {
		return m.processAscii(input)
	}
	return m.processUnicode(input)
}
