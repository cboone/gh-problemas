package utils

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestHexToColor(t *testing.T) {
	tests := []struct {
		input string
		want  lipgloss.Color
	}{
		{"d73a4a", lipgloss.Color("#d73a4a")},
		{"#d73a4a", lipgloss.Color("#d73a4a")},
		{"000000", lipgloss.Color("#000000")},
		{"ffffff", lipgloss.Color("#ffffff")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := HexToColor(tt.input)
			if got != tt.want {
				t.Errorf("HexToColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContrastColor(t *testing.T) {
	tests := []struct {
		name string
		bg   string
		want lipgloss.Color
	}{
		{"white background", "ffffff", lipgloss.Color("#000000")},
		{"black background", "000000", lipgloss.Color("#ffffff")},
		{"saturated red", "ff0000", lipgloss.Color("#000000")},
		{"bright yellow", "ffff00", lipgloss.Color("#000000")},
		{"dark blue", "0000aa", lipgloss.Color("#ffffff")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContrastColor(tt.bg)
			if got != tt.want {
				t.Errorf("ContrastColor(%q) = %v, want %v", tt.bg, got, tt.want)
			}
		})
	}
}

func TestContrastColor_InvalidHex(t *testing.T) {
	got := ContrastColor("xyz")
	if got != lipgloss.Color("#ffffff") {
		t.Errorf("ContrastColor(invalid) = %v, want white", got)
	}
}
