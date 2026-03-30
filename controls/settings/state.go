package settings

type State int

const (
	None State = iota
	Colors
	Characters
	Size
	Adjust
	TextStyle
	Advanced
	Animation
	Alpha
	Theme
	SaveLoad
)

var States = []State{
	Colors,
	Characters,
	Size,
	Adjust,
	TextStyle,
	Advanced,
	Animation,
	Alpha,
}

var stateOrder = []State{Colors, Characters, Size, Adjust, TextStyle, Advanced, Animation, Alpha}

var stateTitles = map[State]string{
	Colors:     "Colors",
	Characters: "Characters",
	Size:       "Size",
	Adjust:     "Adjust",
	TextStyle:  "Text Style",
	Advanced:   "Advanced",
	Animation:  "Animation",
	Alpha:      "Alpha Channel",
	Theme:      "Theme",
	SaveLoad:   "Load / Save Settings",
}
