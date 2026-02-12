package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoConfigFile(t *testing.T) {
	// Point XDG to a temp dir with no config
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if cfg.Defaults.RefreshInterval != 300 {
		t.Errorf("expected refresh_interval 300, got %d", cfg.Defaults.RefreshInterval)
	}
	if cfg.Defaults.PageSize != 50 {
		t.Errorf("expected page_size 50, got %d", cfg.Defaults.PageSize)
	}
	if cfg.Defaults.DateFormat != "relative" {
		t.Errorf("expected date_format relative, got %s", cfg.Defaults.DateFormat)
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected theme dark, got %s", cfg.Theme)
	}
}

func TestLoad_PartialOverride(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "gh-problemas")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `version: 1
defaults:
  page_size: 25
  date_format: "2006-01-02"
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Overridden values
	if cfg.Defaults.PageSize != 25 {
		t.Errorf("expected page_size 25, got %d", cfg.Defaults.PageSize)
	}
	if cfg.Defaults.DateFormat != "2006-01-02" {
		t.Errorf("expected date_format 2006-01-02, got %s", cfg.Defaults.DateFormat)
	}

	// Default values preserved
	if cfg.Defaults.RefreshInterval != 300 {
		t.Errorf("expected refresh_interval 300, got %d", cfg.Defaults.RefreshInterval)
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected theme dark, got %s", cfg.Theme)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "gh-problemas")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("{{invalid yaml"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", tmp)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_XDGConfigHome(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "gh-problemas")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `theme: light
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Theme != "light" {
		t.Errorf("expected theme light, got %s", cfg.Theme)
	}
}
