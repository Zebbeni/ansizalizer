package settings

type State int

const (
	None State = iota
	Palette
	Characters
	Size
	Advanced
)

var States = []State{
	Palette,
	Characters,
	Size,
	Advanced,
}

var stateOrder = []State{Palette, Characters, Size, Advanced}

var stateTitles = map[State]string{
	Palette:    "Palette",
	Characters: "Characters",
	Size:       "Size",
	Advanced:   "Advanced",
}
