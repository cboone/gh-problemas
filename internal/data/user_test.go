package data

import (
	"errors"
	"testing"
)

func TestWhoAmI(t *testing.T) {
	canned := map[string]interface{}{
		"viewer": map[string]interface{}{
			"login": "testuser",
		},
	}

	client := NewUserClient(&mockQuerier{response: canned})
	login, err := client.WhoAmI()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if login != "testuser" {
		t.Errorf("expected testuser, got %s", login)
	}
}

func TestWhoAmI_Error(t *testing.T) {
	client := NewUserClient(&mockQuerier{err: errors.New("auth required")})
	_, err := client.WhoAmI()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
