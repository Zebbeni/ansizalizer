package browser

type State struct {
	DirPath  string
	Filepath string
}

// Callbacks allow the browser to affect global state without direct access

type SetDirPath func(string)
type SetFilePath func(string)

func NewState() State {
	return State{
		DirPath:  "",
		Filepath: "",
	}
}
