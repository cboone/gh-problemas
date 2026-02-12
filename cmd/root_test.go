package cmd

import (
	"encoding/json"
	"errors"
	"testing"
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

func TestResolveRepository_ConfigRepo(t *testing.T) {
	owner, repo, err := resolveRepository("octo/proj", &mockQuerier{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "octo" || repo != "proj" {
		t.Fatalf("unexpected owner/repo: %s/%s", owner, repo)
	}
}

func TestResolveRepository_InvalidConfigRepo(t *testing.T) {
	_, _, err := resolveRepository("not-valid", &mockQuerier{})
	if err == nil {
		t.Fatal("expected error for invalid defaults.repo")
	}
}

func TestResolveRepository_AtMe(t *testing.T) {
	q := &mockQuerier{response: map[string]interface{}{
		"viewer": map[string]interface{}{
			"login": "alice",
		},
	}}

	owner, repo, err := resolveRepository("@me/proj", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "alice" || repo != "proj" {
		t.Fatalf("unexpected owner/repo: %s/%s", owner, repo)
	}
}

func TestResolveRepository_AtMeError(t *testing.T) {
	_, _, err := resolveRepository("@me/proj", &mockQuerier{err: errors.New("auth required")})
	if err == nil {
		t.Fatal("expected error for @me resolution failure")
	}
}
