package utils

import "github.com/charmbracelet/glamour"

// RenderMarkdown renders a markdown string for terminal display.
func RenderMarkdown(content string, width int) (string, error) {
	if content == "" {
		return "", nil
	}

	if width < 1 {
		width = 1
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	return r.Render(content)
}
