package settings

type State int

const (
	None State = iota
	Colors
	Characters
	Size
	Advanced
)

var States = []State{
	Colors,
	Characters,
	Size,
	Advanced,
}

var stateOrder = []State{Colors, Characters, Size, Advanced}

var stateTitles = map[State]string{
	Colors:     "Colors",
	Characters: "Characters",
	Size:       "Size",
	Advanced:   "Advanced",
}
