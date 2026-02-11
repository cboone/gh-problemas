package utils

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"30 seconds", 30 * time.Second, "30s ago"},
		{"90 minutes", 90 * time.Minute, "1h ago"},
		{"25 hours", 25 * time.Hour, "1d ago"},
		{"60 days", 60 * 24 * time.Hour, "2mo ago"},
		{"400 days", 400 * 24 * time.Hour, "1y ago"},
		{"5 minutes", 5 * time.Minute, "5m ago"},
		{"3 hours", 3 * time.Hour, "3h ago"},
		{"10 days", 10 * 24 * time.Hour, "10d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RelativeTime(time.Now().Add(-tt.duration))
			if got != tt.want {
				t.Errorf("RelativeTime(%v ago) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestRelativeTime_Zero(t *testing.T) {
	got := RelativeTime(time.Time{})
	if got != "" {
		t.Errorf("RelativeTime(zero) = %q, want empty string", got)
	}
}

func TestRelativeTime_Future(t *testing.T) {
	got := RelativeTime(time.Now().Add(1 * time.Hour))
	if got != "just now" {
		t.Errorf("RelativeTime(future) = %q, want 'just now'", got)
	}
}
