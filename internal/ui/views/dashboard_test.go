package views

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cboone/gh-problemas/internal/data"
	"github.com/cboone/gh-problemas/internal/ui"
)

type mockQuerier struct {
	response interface{}
	err      error
}

func (m *mockQuerier) Do(_ string, _ map[string]interface{}, resp interface{}) error {
	if m.err != nil {
		return m.err
	}
	b, err := json.Marshal(m.response)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, resp)
}

func TestDashboard_IssuesLoadedMsg_PopulatesList(t *testing.T) {
	client := data.NewIssueClient(&mockQuerier{}, "owner", "repo")
	styles := ui.DefaultStyles()
	keys := ui.DefaultKeyMap()
	dv := NewDashboardView(client, styles, keys, 80, 24)

	msg := ui.IssuesLoadedMsg{
		Result: data.IssueListResult{
			Issues: []data.Issue{
				{Number: 1, Title: "First", Author: "alice", CreatedAt: time.Now()},
				{Number: 2, Title: "Second", Author: "bob", CreatedAt: time.Now()},
				{Number: 3, Title: "Third", Author: "carol", CreatedAt: time.Now()},
			},
		},
	}

	updated, _ := dv.Update(msg)
	d := updated.(*DashboardView)

	if d.loading {
		t.Error("expected loading to be false after IssuesLoadedMsg")
	}
	if d.errMsg != "" {
		t.Errorf("expected no error, got %q", d.errMsg)
	}
	items := d.list.Items()
	if len(items) != 3 {
		t.Fatalf("expected 3 items in list, got %d", len(items))
	}
}

func TestDashboard_IssuesLoadedMsg_Error(t *testing.T) {
	client := data.NewIssueClient(&mockQuerier{}, "owner", "repo")
	styles := ui.DefaultStyles()
	keys := ui.DefaultKeyMap()
	dv := NewDashboardView(client, styles, keys, 80, 24)

	msg := ui.IssuesLoadedMsg{
		Err: errors.New("api error"),
	}

	updated, _ := dv.Update(msg)
	d := updated.(*DashboardView)

	if d.errMsg == "" {
		t.Error("expected error message to be set")
	}
}

func TestDashboard_EnterProducesNavigateMsg(t *testing.T) {
	client := data.NewIssueClient(&mockQuerier{}, "owner", "repo")
	styles := ui.DefaultStyles()
	keys := ui.DefaultKeyMap()
	dv := NewDashboardView(client, styles, keys, 80, 24)

	// First load issues
	loadMsg := ui.IssuesLoadedMsg{
		Result: data.IssueListResult{
			Issues: []data.Issue{
				{Number: 42, Title: "Test issue", Author: "alice", CreatedAt: time.Now()},
			},
		},
	}
	dv.Update(loadMsg)

	// Press enter
	_, cmd := dv.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command from Enter key, got nil")
	}
	msg := cmd()
	nav, ok := msg.(ui.NavigateToDetailMsg)
	if !ok {
		t.Fatalf("expected NavigateToDetailMsg, got %T", msg)
	}
	if nav.IssueNumber != 42 {
		t.Errorf("expected issue number 42, got %d", nav.IssueNumber)
	}
}

func TestDashboard_RefreshTriggersLoad(t *testing.T) {
	client := data.NewIssueClient(&mockQuerier{}, "owner", "repo")
	styles := ui.DefaultStyles()
	keys := ui.DefaultKeyMap()
	dv := NewDashboardView(client, styles, keys, 80, 24)

	// Load initial issues then refresh
	loadMsg := ui.IssuesLoadedMsg{
		Result: data.IssueListResult{
			Issues: []data.Issue{
				{Number: 1, Title: "Test", Author: "alice", CreatedAt: time.Now()},
			},
		},
	}
	dv.Update(loadMsg)

	// Press R to refresh
	updated, cmd := dv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
	d := updated.(*DashboardView)

	if !d.loading {
		t.Error("expected loading to be true after refresh")
	}
	if cmd == nil {
		t.Error("expected command from refresh")
	}
}
