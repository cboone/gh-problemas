package utils

import (
	"strings"
	"testing"
)

func TestRenderMarkdown_Simple(t *testing.T) {
	out, err := RenderMarkdown("**bold** text", 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !strings.Contains(out, "bold") {
		t.Errorf("expected output to contain 'bold', got: %q", out)
	}
}

func TestRenderMarkdown_Empty(t *testing.T) {
	out, err := RenderMarkdown("", 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output for empty input, got: %q", out)
	}
}

func TestRenderMarkdown_WidthWrapping(t *testing.T) {
	long := "This is a very long line of text that should be wrapped when rendered with a narrow width setting to test word wrapping behavior."
	out, err := RenderMarkdown(long, 40)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Errorf("expected text to wrap into multiple lines at width 40, got %d lines", len(lines))
	}
}
