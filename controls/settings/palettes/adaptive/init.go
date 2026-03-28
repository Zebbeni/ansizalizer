package adaptive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

func newInput(state State) textinput.Model {
	input := textinput.New()
	input.Prompt = stateNames[state]
	input.Cursor.Blink = true
	input.CharLimit = 3
	input.SetValue(fmt.Sprintf("16"))
	return input
}
