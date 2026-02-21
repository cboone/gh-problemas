package data

import (
	"encoding/json"
	"errors"
	"testing"
)

// mockQuerier returns canned responses for testing.
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

func TestList_ThreeIssues(t *testing.T) {
	canned := map[string]interface{}{
		"repository": map[string]interface{}{
			"issues": map[string]interface{}{
				"pageInfo": map[string]interface{}{
					"hasNextPage": true,
					"endCursor":   "cursor123",
				},
				"nodes": []map[string]interface{}{
					{
						"number": 1, "title": "First issue", "state": "OPEN",
						"createdAt": "2025-01-01T00:00:00Z", "updatedAt": "2025-01-02T00:00:00Z",
						"author":    map[string]string{"login": "alice"},
						"labels":    map[string]interface{}{"nodes": []interface{}{}},
						"assignees": map[string]interface{}{"nodes": []interface{}{}},
						"comments":  map[string]int{"totalCount": 3},
						"reactions": map[string]int{"totalCount": 5},
					},
					{
						"number": 2, "title": "Second issue", "state": "OPEN",
						"createdAt": "2025-01-03T00:00:00Z", "updatedAt": "2025-01-04T00:00:00Z",
						"author": map[string]string{"login": "bob"},
						"labels": map[string]interface{}{
							"nodes": []map[string]string{{"name": "bug", "color": "d73a4a"}},
						},
						"assignees": map[string]interface{}{
							"nodes": []map[string]string{{"login": "bob"}},
						},
						"milestone": map[string]string{"title": "v1.0"},
						"comments":  map[string]int{"totalCount": 1},
						"reactions": map[string]int{"totalCount": 0},
					},
					{
						"number": 3, "title": "Third issue", "state": "CLOSED",
						"createdAt": "2025-01-05T00:00:00Z", "updatedAt": "2025-01-06T00:00:00Z",
						"author":    map[string]string{"login": "carol"},
						"labels":    map[string]interface{}{"nodes": []interface{}{}},
						"assignees": map[string]interface{}{"nodes": []interface{}{}},
						"comments":  map[string]int{"totalCount": 0},
						"reactions": map[string]int{"totalCount": 2},
					},
				},
			},
		},
	}

	client := NewIssueClient(&mockQuerier{response: canned}, "owner", "repo")
	result, err := client.List(IssueListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Issues) != 3 {
		t.Fatalf("expected 3 issues, got %d", len(result.Issues))
	}
	if !result.PageInfo.HasNextPage {
		t.Error("expected HasNextPage to be true")
	}
	if result.PageInfo.EndCursor != "cursor123" {
		t.Errorf("expected EndCursor cursor123, got %s", result.PageInfo.EndCursor)
	}

	// Verify first issue
	first := result.Issues[0]
	if first.Number != 1 || first.Title != "First issue" || first.Author != "alice" {
		t.Errorf("unexpected first issue: %+v", first)
	}
	if first.CommentCount != 3 || first.ReactionCount != 5 {
		t.Errorf("unexpected counts: comments=%d reactions=%d", first.CommentCount, first.ReactionCount)
	}

	// Verify second issue with labels, assignees, milestone
	second := result.Issues[1]
	if len(second.Labels) != 1 || second.Labels[0].Name != "bug" {
		t.Errorf("unexpected labels: %+v", second.Labels)
	}
	if len(second.Assignees) != 1 || second.Assignees[0] != "bob" {
		t.Errorf("unexpected assignees: %+v", second.Assignees)
	}
	if second.Milestone != "v1.0" {
		t.Errorf("expected milestone v1.0, got %s", second.Milestone)
	}
}

func TestList_Empty(t *testing.T) {
	canned := map[string]interface{}{
		"repository": map[string]interface{}{
			"issues": map[string]interface{}{
				"pageInfo": map[string]interface{}{"hasNextPage": false, "endCursor": ""},
				"nodes":    []interface{}{},
			},
		},
	}

	client := NewIssueClient(&mockQuerier{response: canned}, "owner", "repo")
	result, err := client.List(IssueListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(result.Issues))
	}
	if result.PageInfo.HasNextPage {
		t.Error("expected HasNextPage to be false")
	}
}

func TestList_GraphQLError(t *testing.T) {
	client := NewIssueClient(&mockQuerier{err: errors.New("graphql: auth required")}, "owner", "repo")
	_, err := client.List(IssueListOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "graphql: auth required" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGet_FullBody(t *testing.T) {
	canned := map[string]interface{}{
		"repository": map[string]interface{}{
			"issue": map[string]interface{}{
				"number": 42, "title": "Important bug", "state": "OPEN",
				"createdAt": "2025-03-01T10:00:00Z", "updatedAt": "2025-03-02T15:30:00Z",
				"author": map[string]string{"login": "dave"},
				"labels": map[string]interface{}{
					"nodes": []map[string]string{
						{"name": "bug", "color": "d73a4a"},
						{"name": "urgent", "color": "e4e669"},
					},
				},
				"assignees": map[string]interface{}{
					"nodes": []map[string]string{{"login": "dave"}, {"login": "eve"}},
				},
				"milestone": map[string]string{"title": "v2.0"},
				"comments":  map[string]int{"totalCount": 7},
				"reactions": map[string]int{"totalCount": 12},
				"body":      "## Steps to reproduce\n\n1. Do this\n2. Do that\n\n**Expected:** works\n**Actual:** broken",
			},
		},
	}

	client := NewIssueClient(&mockQuerier{response: canned}, "owner", "repo")
	issue, err := client.Get(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if issue.Number != 42 {
		t.Errorf("expected number 42, got %d", issue.Number)
	}
	if issue.Title != "Important bug" {
		t.Errorf("expected title 'Important bug', got %s", issue.Title)
	}
	if issue.Author != "dave" {
		t.Errorf("expected author dave, got %s", issue.Author)
	}
	if len(issue.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(issue.Labels))
	}
	if len(issue.Assignees) != 2 {
		t.Errorf("expected 2 assignees, got %d", len(issue.Assignees))
	}
	if issue.Milestone != "v2.0" {
		t.Errorf("expected milestone v2.0, got %s", issue.Milestone)
	}
	if issue.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestGet_Error(t *testing.T) {
	client := NewIssueClient(&mockQuerier{err: errors.New("graphql: not found")}, "owner", "repo")
	_, err := client.Get(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
