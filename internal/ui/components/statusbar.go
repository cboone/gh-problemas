package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar renders a bottom bar with repo context, key hints, and messages.
type StatusBar struct {
	repoName string
	keyHints []string
	message  string
	width    int
	style    lipgloss.Style
}

// NewStatusBar creates a new status bar.
func NewStatusBar(style lipgloss.Style) *StatusBar {
	return &StatusBar{style: style}
}

// SetRepoName sets the repository name displayed on the left.
func (s *StatusBar) SetRepoName(name string) {
	s.repoName = name
}

// SetKeyHints sets the key hints displayed in the center.
func (s *StatusBar) SetKeyHints(hints []string) {
	s.keyHints = hints
}

// SetMessage sets a transient message (error/success) on the right.
func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
}

// SetWidth sets the status bar width.
func (s *StatusBar) SetWidth(w int) {
	s.width = w
}

// View renders the status bar.
func (s *StatusBar) View() string {
	left := s.repoName
	center := strings.Join(s.keyHints, " | ")
	right := s.message

	// Calculate available space
	available := s.width - lipgloss.Width(left) - lipgloss.Width(right)
	if available < 0 {
		available = 0
	}

	// Pad center to fill available space
	centerPad := ""
	if available > lipgloss.Width(center) {
		padding := available - lipgloss.Width(center)
		leftPad := padding / 2
		rightPad := padding - leftPad
		centerPad = strings.Repeat(" ", leftPad) + center + strings.Repeat(" ", rightPad)
	} else {
		centerPad = center
	}

	bar := left + centerPad + right
	return s.style.Width(s.width).Render(bar)
}
