package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/starkross/augur/internal/config"
)

func TestLoadYAML_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte("key: value\nnested:\n  a: 1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	m, err := config.LoadYAML(path)
	if err != nil {
		t.Fatalf("LoadYAML() error: %v", err)
	}

	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestLoadYAML_NotFound(t *testing.T) {
	_, err := config.LoadYAML("/nonexistent/file.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadYAML_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte(":\n  :\n    [invalid"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := config.LoadYAML(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadYAML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(path, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := config.LoadYAML(path)
	if err == nil {
		t.Error("expected error for empty YAML")
	}
}
