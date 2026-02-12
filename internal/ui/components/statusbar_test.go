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

func TestStatusBar_SeparatorWhenCenterSqueezed(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetRepoName("owner/repo")
	sb.SetKeyHints([]string{"j/k: navigate", "enter: open"})
	sb.SetInfo("some status message here")
	// Width tight enough that left+right fill the bar, squeezing center out
	sb.SetWidth(lipgloss.Width("owner/repo") + lipgloss.Width("info: some status message here"))

	v := sb.View()
	// Left and right should not run together without separation
	if strings.Contains(v, "owner/repoinfo:") || strings.Contains(v, "po…info:") {
		t.Fatalf("expected separator between left and right, got %q", v)
	}
}

func TestStatusBar_SeparatorWhenCenterTruncated(t *testing.T) {
	sb := NewStatusBar(lipgloss.NewStyle())
	sb.SetRepoName("owner/repo")
	sb.SetKeyHints([]string{"j/k: navigate", "enter: open"})
	sb.SetInfo("api: msg")
	// Width allows some center content but requires truncation
	sb.SetWidth(lipgloss.Width("owner/repo") + lipgloss.Width("info: api: msg") + 8)

	v := sb.View()
	left := "owner/repo"
	right := "info: api: msg"
	// Segments should not be directly adjacent to truncated center
	if strings.Contains(v, left+"j") || strings.Contains(v, left+"…") {
		t.Fatalf("expected space separator between left and truncated center, got %q", v)
	}
	if idx := strings.Index(v, right); idx > 0 {
		before := v[idx-1 : idx]
		if before != " " {
			t.Fatalf("expected space before right segment, got %q at position %d in %q", before, idx-1, v)
		}
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
