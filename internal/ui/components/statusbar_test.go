package components

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStatusBar_SetError_Network(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetError(errors.New("dial tcp: i/o timeout"))
	sb.SetWidth(120)

	v := sb.View()
	if !strings.Contains(v, "network: dial tcp") {
		t.Fatalf("expected network classification, got %q", v)
	}
}

func TestStatusBar_SetError_API401(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetError(errors.New("HTTP 401 unauthorized"))
	sb.SetWidth(120)

	v := sb.View()
	if !strings.Contains(v, "api: Run gh auth login") {
		t.Fatalf("expected auth guidance, got %q", v)
	}
}

func TestStatusBar_SetError_API404(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetError(errors.New("HTTP 404 not found"))
	sb.SetWidth(120)

	v := sb.View()
	if !strings.Contains(v, "api: Repository not found") {
		t.Fatalf("expected repo not found guidance, got %q", v)
	}
}

func TestStatusBar_TruncatesRightMessage(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetRepoName("owner/repository-with-very-long-name")
	sb.SetKeyHints([]string{"j/k: navigate", "enter: open"})
	sb.SetInfo("this is a very long status message that should be truncated for narrow widths")
	sb.SetWidth(40)

	v := sb.View()
	if lipgloss.Width(v) == 0 {
		t.Fatal("expected rendered status bar")
	}
	if !strings.Contains(v, "…") {
		t.Fatalf("expected ellipsis truncation, got %q", v)
	}
}
