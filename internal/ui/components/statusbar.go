package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar renders a bottom bar with repo context, key hints, and messages.
type StatusBar struct {
	repoName      string
	keyHints      []string
	message       string
	messagePrefix string
	width         int
	style         lipgloss.Style
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
	s.messagePrefix = ""
}

// SetInfo sets an informational status message.
func (s *StatusBar) SetInfo(msg string) {
	s.message = msg
	s.messagePrefix = "info"
}

// SetLoading sets a loading status message.
func (s *StatusBar) SetLoading(msg string) {
	s.message = msg
	s.messagePrefix = "loading"
}

// SetError sets an error message with API/network classification.
func (s *StatusBar) SetError(err error) {
	if err == nil {
		s.SetMessage("")
		return
	}

	text := strings.TrimSpace(err.Error())
	if text == "" {
		s.SetMessage("")
		return
	}

	if isNetworkError(text) {
		s.messagePrefix = "network"
		s.message = text
		return
	}

	if strings.Contains(text, "401") {
		s.messagePrefix = "api"
		s.message = "Run gh auth login to re-authenticate"
		return
	}

	if strings.Contains(text, "403") {
		s.messagePrefix = "api"
		s.message = "Check your permissions for this repository"
		return
	}

	if strings.Contains(text, "404") {
		s.messagePrefix = "api"
		s.message = "Repository not found"
		return
	}

	s.messagePrefix = "api"
	s.message = text
}

// SetWidth sets the status bar width.
func (s *StatusBar) SetWidth(w int) {
	s.width = w
}

// View renders the status bar.
func (s *StatusBar) View() string {
	left := s.repoName
	center := strings.Join(s.keyHints, " | ")
	right := s.renderMessage()

	if s.width <= 0 {
		bar := strings.TrimSpace(left + " " + center + " " + right)
		return s.style.Render(bar)
	}

	maxLeft := s.width / 4
	if maxLeft < 12 {
		maxLeft = 12
	}
	left = truncateText(left, maxLeft)

	maxRight := s.width / 3
	if maxRight < 20 {
		maxRight = 20
	}
	right = truncateText(right, maxRight)

	// Calculate available space
	available := s.width - lipgloss.Width(left) - lipgloss.Width(right)
	if available < 0 {
		left = truncateText(left, s.width-lipgloss.Width(right))
		available = s.width - lipgloss.Width(left) - lipgloss.Width(right)
		if available < 0 {
			available = 0
		}
	}

	// Ensure minimum separator between left and right when center is fully squeezed
	if available == 0 && left != "" && right != "" && lipgloss.Width(left) > 1 {
		left = truncateText(left, lipgloss.Width(left)-1)
		available = 1
	}

	// Pad center to fill available space
	centerPad := ""
	if available > lipgloss.Width(center) {
		padding := available - lipgloss.Width(center)
		leftPad := padding / 2
		rightPad := padding - leftPad
		centerPad = strings.Repeat(" ", leftPad) + center + strings.Repeat(" ", rightPad)
	} else {
		centerPad = truncateText(center, available)
	}

	bar := left + centerPad + right
	return s.style.Width(s.width).Render(bar)
}

func (s *StatusBar) renderMessage() string {
	if s.message == "" {
		return ""
	}

	switch s.messagePrefix {
	case "loading":
		return "loading: " + s.message
	case "network":
		return "network: " + s.message
	case "api":
		return "api: " + s.message
	default:
		return s.message
	}
}

func truncateText(text string, max int) string {
	if max <= 0 {
		return ""
	}
	if lipgloss.Width(text) <= max {
		return text
	}
	if max == 1 {
		return "…"
	}
	runes := []rune(text)
	if len(runes) >= max {
		return string(runes[:max-1]) + "…"
	}
	return text
}

func isNetworkError(errText string) bool {
	text := strings.ToLower(errText)
	keywords := []string{
		"dial tcp",
		"connection refused",
		"no such host",
		"i/o timeout",
		"network is unreachable",
		"timeout",
		"tls",
		"temporary failure",
	}

	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}
