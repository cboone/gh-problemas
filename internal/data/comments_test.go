package data

import (
	"errors"
	"testing"
)

func TestCommentList_ThreeComments(t *testing.T) {
	canned := map[string]interface{}{
		"repository": map[string]interface{}{
			"issue": map[string]interface{}{
				"comments": map[string]interface{}{
					"pageInfo": map[string]interface{}{"hasNextPage": false, "endCursor": ""},
					"nodes": []map[string]interface{}{
						{
							"author":    map[string]string{"login": "alice"},
							"body":      "First comment",
							"createdAt": "2025-01-01T00:00:00Z",
							"updatedAt": "2025-01-01T00:00:00Z",
							"reactions": map[string]int{"totalCount": 2},
						},
						{
							"author":    map[string]string{"login": "bob"},
							"body":      "Second comment",
							"createdAt": "2025-01-02T00:00:00Z",
							"updatedAt": "2025-01-02T00:00:00Z",
							"reactions": map[string]int{"totalCount": 0},
						},
						{
							"author":    map[string]string{"login": "carol"},
							"body":      "Third comment with **markdown**",
							"createdAt": "2025-01-03T00:00:00Z",
							"updatedAt": "2025-01-03T00:00:00Z",
							"reactions": map[string]int{"totalCount": 5},
						},
					},
				},
			},
		},
	}

	client := NewCommentClient(&mockQuerier{response: canned}, "owner", "repo")
	result, err := client.List(1, 25, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Comments) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(result.Comments))
	}
	if result.Comments[0].Author != "alice" {
		t.Errorf("expected first author alice, got %s", result.Comments[0].Author)
	}
	if result.Comments[2].Reactions != 5 {
		t.Errorf("expected 5 reactions on third comment, got %d", result.Comments[2].Reactions)
	}
}

func TestCommentList_Empty(t *testing.T) {
	canned := map[string]interface{}{
		"repository": map[string]interface{}{
			"issue": map[string]interface{}{
				"comments": map[string]interface{}{
					"pageInfo": map[string]interface{}{"hasNextPage": false, "endCursor": ""},
					"nodes":    []interface{}{},
				},
			},
		},
	}

	client := NewCommentClient(&mockQuerier{response: canned}, "owner", "repo")
	result, err := client.List(1, 25, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Comments) != 0 {
		t.Fatalf("expected 0 comments, got %d", len(result.Comments))
	}
}

func TestCommentList_Error(t *testing.T) {
	client := NewCommentClient(&mockQuerier{err: errors.New("graphql error")}, "owner", "repo")
	_, err := client.List(1, 25, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
