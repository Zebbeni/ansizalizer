package adaptive

import (
	"github.com/Zebbeni/ansizalizer/controls/numberinput"
	"github.com/Zebbeni/ansizalizer/style"
)

func newInput(state State) numberinput.Model {
	m := numberinput.New(numberinput.Options{
		Prompt:    stateNames[state] + " ",
		CharLimit: 3,
		IsFloat:   false,
		Min:       numberinput.FloatPtr(1),
		Default:   16,
	})
	styles := m.Styles()
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	m.SetStyles(styles)
	return m
}
