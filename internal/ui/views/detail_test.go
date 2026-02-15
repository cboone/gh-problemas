package views

import (
	"strings"
	"testing"
	"time"

	"github.com/cboone/gh-problemas/internal/data"
	"github.com/cboone/gh-problemas/internal/ui"
)

func TestDetailView_UsesConfiguredDateFormat(t *testing.T) {
	styles := ui.DefaultStyles()
	keys := ui.DefaultKeyMap()
	dv := NewDetailViewWithCommentsAndDateFormat(nil, nil, styles, keys, 1, 100, 30, "2006-01-02")

	dv.issue = &data.Issue{
		Number:    12,
		Title:     "Date format test",
		State:     "OPEN",
		Author:    "alice",
		CreatedAt: time.Date(2025, time.March, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2025, time.March, 2, 12, 0, 0, 0, time.UTC),
		Body:      "hello",
	}

	dv.renderContent()
	out := dv.viewport.View()
	if !strings.Contains(out, "Created: 2025-03-01") {
		t.Fatalf("expected created date in custom format, got: %q", out)
	}
	if !strings.Contains(out, "Updated: 2025-03-02") {
		t.Fatalf("expected updated date in custom format, got: %q", out)
	}
}
