package settings

type State int

const (
	None State = iota
	Palette
	Characters
	Size
	Sampling
)

var States = []State{
	Palette,
	Characters,
	Size,
	Sampling,
}

var stateOrder = []State{Palette, Characters, Size, Sampling}

var stateTitles = map[State]string{
	Palette:    "Palette",
	Characters: "Characters",
	Size:       "Size",
	Sampling:   "Sampling",
}
