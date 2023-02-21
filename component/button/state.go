package button

type State int

const (
	IsInactive State = iota
	IsHighlighted
	IsSelected
)
