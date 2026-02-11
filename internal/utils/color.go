package utils

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HexToColor converts a hex color string (with or without '#') to a lipgloss.Color.
func HexToColor(hex string) lipgloss.Color {
	hex = strings.TrimPrefix(hex, "#")
	return lipgloss.Color("#" + hex)
}

// ContrastColor returns black or white depending on the background luminance,
// using the W3C relative luminance algorithm.
func ContrastColor(backgroundHex string) lipgloss.Color {
	backgroundHex = strings.TrimPrefix(backgroundHex, "#")
	if len(backgroundHex) != 6 {
		return lipgloss.Color("#ffffff")
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(backgroundHex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return lipgloss.Color("#ffffff")
	}

	luminance := relativeLuminance(r, g, b)
	if luminance > 0.179 {
		return lipgloss.Color("#000000")
	}
	return lipgloss.Color("#ffffff")
}

func relativeLuminance(r, g, b uint8) float64 {
	rs := linearize(float64(r) / 255.0)
	gs := linearize(float64(g) / 255.0)
	bs := linearize(float64(b) / 255.0)
	return 0.2126*rs + 0.7152*gs + 0.0722*bs
}

func linearize(c float64) float64 {
	if c <= 0.03928 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}
