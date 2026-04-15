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
	if err := os.WriteFile(path, []byte("key: value\nnested:\n  a: 1\n"), 0o600); err != nil {
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
	if err := os.WriteFile(path, []byte(":\n  :\n    [invalid"), 0o600); err != nil {
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
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := config.LoadYAML(path)
	if err == nil {
		t.Error("expected error for empty YAML")
	}
}

func writeFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadMerged_MapsDeepMerge(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "receivers:\n  otlp:\n    protocols:\n      grpc: {}\n")
	b := writeFile(t, dir, "b.yaml", "receivers:\n  kafka:\n    brokers: [k1]\n")

	m, err := config.LoadMerged([]string{a, b})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}

	recv, ok := m["receivers"].(map[string]any)
	if !ok {
		t.Fatalf("receivers not a map: %T", m["receivers"])
	}
	if _, ok := recv["otlp"]; !ok {
		t.Error("expected otlp preserved from a.yaml")
	}
	if _, ok := recv["kafka"]; !ok {
		t.Error("expected kafka added from b.yaml")
	}
}

func TestLoadMerged_ScalarLaterWins(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "service:\n  telemetry:\n    logs:\n      level: info\n")
	b := writeFile(t, dir, "b.yaml", "service:\n  telemetry:\n    logs:\n      level: debug\n")

	m, err := config.LoadMerged([]string{a, b})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}
	service, _ := m["service"].(map[string]any)
	telemetry, _ := service["telemetry"].(map[string]any)
	logs, _ := telemetry["logs"].(map[string]any)
	if got := logs["level"]; got != "debug" {
		t.Errorf("expected level=debug, got %v", got)
	}
}

func TestLoadMerged_SliceReplace(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "service:\n  pipelines:\n    traces:\n      receivers: [otlp]\n      processors: [batch]\n")
	b := writeFile(t, dir, "b.yaml", "service:\n  pipelines:\n    traces:\n      receivers: [kafka]\n")

	m, err := config.LoadMerged([]string{a, b})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}
	service, _ := m["service"].(map[string]any)
	pipelines, _ := service["pipelines"].(map[string]any)
	traces, _ := pipelines["traces"].(map[string]any)
	recv, _ := traces["receivers"].([]any)
	if len(recv) != 1 || recv[0] != "kafka" {
		t.Errorf("expected receivers=[kafka], got %v", recv)
	}
	procs, _ := traces["processors"].([]any)
	if len(procs) != 1 || procs[0] != "batch" {
		t.Errorf("expected processors=[batch] preserved, got %v", procs)
	}
}

func TestLoadMerged_SingleFile(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "key: value\n")
	m, err := config.LoadMerged([]string{a})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestLoadMerged_EmptyList(t *testing.T) {
	if _, err := config.LoadMerged(nil); err == nil {
		t.Error("expected error for empty path list")
	}
}

func TestLoadMerged_TypeMismatchLaterWins(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "foo:\n  nested: true\n")
	b := writeFile(t, dir, "b.yaml", "foo: 42\n")

	m, err := config.LoadMerged([]string{a, b})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}
	if m["foo"] != float64(42) {
		t.Errorf("expected foo=42, got %v (%T)", m["foo"], m["foo"])
	}
}

func TestLoadMerged_ThreeFileChain(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "x: 1\n")
	b := writeFile(t, dir, "b.yaml", "x: 2\n")
	c := writeFile(t, dir, "c.yaml", "x: 3\n")

	m, err := config.LoadMerged([]string{a, b, c})
	if err != nil {
		t.Fatalf("LoadMerged: %v", err)
	}
	if m["x"] != float64(3) {
		t.Errorf("expected x=3 (last wins), got %v", m["x"])
	}
}

func TestLoadMerged_MissingFile(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.yaml", "x: 1\n")
	if _, err := config.LoadMerged([]string{a, "/nonexistent.yaml"}); err == nil {
		t.Error("expected error when a later file is missing")
	}
}
