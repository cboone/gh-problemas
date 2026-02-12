package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds all application key bindings.
type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Open      key.Binding
	Back      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
	Refresh   key.Binding
	Help      key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	GoToTop   key.Binding
	GoToEnd   key.Binding
	NextPage  key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:        key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/up", "up")),
		Down:      key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/down", "down")),
		Open:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Back:      key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "back")),
		Quit:      key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "force quit")),
		Refresh:   key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "refresh")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		PageUp:    key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup", "page up")),
		PageDown:  key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn", "page down")),
		GoToTop:   key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "go to top")),
		GoToEnd:   key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "go to end")),
		NextPage:  key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "load more")),
	}
}
