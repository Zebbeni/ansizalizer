package loader

import (
	"github.com/Zebbeni/ansizalizer/palette"
)

type item struct {
	palette palette.Model
}

func (i item) FilterValue() string {
	return i.palette.Name()
}

func (i item) Title() string {
	return i.palette.Name()
}

func (i item) Description() string {
	return i.palette.View()
}
