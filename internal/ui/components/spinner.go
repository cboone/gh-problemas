package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner wraps bubbles/spinner with a label and active state.
type Spinner struct {
	model  spinner.Model
	label  string
	active bool
}

// NewSpinner creates a new spinner with the given style.
func NewSpinner(style lipgloss.Style) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = style
	return &Spinner{model: s}
}

// Start activates the spinner with the given label.
func (s *Spinner) Start(label string) tea.Cmd {
	s.active = true
	s.label = label
	return s.model.Tick
}

// Stop deactivates the spinner.
func (s *Spinner) Stop() {
	s.active = false
	s.label = ""
}

// IsActive returns whether the spinner is currently active.
func (s *Spinner) IsActive() bool {
	return s.active
}

// Update processes spinner tick messages.
func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	if !s.active {
		return nil
	}
	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return cmd
}

// View renders the spinner with its label.
func (s *Spinner) View() string {
	if !s.active {
		return ""
	}
	return s.model.View() + " " + s.label
}
