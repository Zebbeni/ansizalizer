package keyboard

import tea "github.com/charmbracelet/bubbletea"

type Handler interface {
	// HandleKeyMsg should return true if the Handler (or any of its child
	// Handlers) handled the given keypress
	HandleKeyMsg(msg tea.KeyMsg) bool
}
