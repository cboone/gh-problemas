package utils

import (
	"fmt"
	"time"
)

const defaultDateFormat = "relative"

// RelativeTime returns a human-readable relative time string.
func RelativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	d := time.Since(t)
	if d < 0 {
		return "just now"
	}

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}

// FormatTime returns a timestamp formatted according to the configured format.
// Supported modes:
// - "relative" (or empty): human-readable relative time
// - any Go time layout string (for example "2006-01-02")
func FormatTime(t time.Time, format string) string {
	if t.IsZero() {
		return ""
	}

	if format == "" || format == defaultDateFormat {
		return RelativeTime(t)
	}

	return t.Local().Format(format)
}
