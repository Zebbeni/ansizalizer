package adaptive

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"

	"github.com/Zebbeni/ansizalizer/style"
)

func newInput(state State) textinput.Model {
	input := textinput.New()
	input.Prompt = stateNames[state]
	styles := input.Styles()
	styles.Cursor.Blink = true
	styles.Cursor.Color = style.SelectedColor1
	input.SetStyles(styles)
	input.CharLimit = 3
	input.SetValue(fmt.Sprintf("16"))
	return input
}
