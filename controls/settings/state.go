package settings

type State int

const (
	None State = iota
	Colors
	Characters
	Size
	Advanced
	Alpha
)

var States = []State{
	Colors,
	Characters,
	Size,
	Advanced,
	Alpha,
}

var stateOrder = []State{Colors, Characters, Size, Advanced, Alpha}

var stateTitles = map[State]string{
	Colors:     "Colors",
	Characters: "Characters",
	Size:       "Size",
	Advanced:   "Advanced",
	Alpha:      "Alpha Channel",
}
