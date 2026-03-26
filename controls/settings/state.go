package settings

type State int

const (
	None State = iota
	Colors
	Characters
	Size
	Adjust
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
	Advanced,
	Animation,
	Alpha,
	SaveLoad,
}

var stateOrder = []State{Colors, Characters, Size, Adjust, Advanced, Animation, Alpha, SaveLoad}

var stateTitles = map[State]string{
	Colors:     "Colors",
	Characters: "Characters",
	Size:       "Size",
	Adjust:     "Adjust",
	Advanced:   "Advanced",
	Animation:  "Animation",
	Alpha:      "Alpha Channel",
	Theme:      "Theme",
	SaveLoad:   "Load / Save Settings",
}
