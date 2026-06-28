package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/starkross/augur/internal/config"
)

func TestLoadEnvFile_Basic(t *testing.T) {
	dir := t.TempDir()
	body := "# comment\n" +
		"\n" +
		"FOO=bar\n" +
		"QUOTED=\"value with spaces\"\n" +
		"SINGLE='single quoted'\n" +
		"export EXPORTED=ok\n" +
		"EMPTY=\n"
	p := writeFile(t, dir, ".env", body)

	m, err := config.LoadEnvFile(p)
	if err != nil {
		t.Fatalf("LoadEnvFile: %v", err)
	}
	want := map[string]string{
		"FOO":      "bar",
		"QUOTED":   "value with spaces",
		"SINGLE":   "single quoted",
		"EXPORTED": "ok",
		"EMPTY":    "",
	}
	for k, v := range want {
		if got := m[k]; got != v {
			t.Errorf("key %q: want %q, got %q", k, v, got)
		}
	}
}

func TestLoadEnvFile_InvalidLine(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "no_equals_sign\n")
	if _, err := config.LoadEnvFile(p); err == nil {
		t.Error("expected error for malformed line")
	}
}

func TestLoadEnvFile_InvalidKey(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "1BAD=value\n")
	if _, err := config.LoadEnvFile(p); err == nil {
		t.Error("expected error for invalid key")
	}
}

func TestLoadEnvFiles_LaterOverrides(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a.env", "X=1\nY=2\n")
	b := writeFile(t, dir, "b.env", "X=overridden\n")
	m, err := config.LoadEnvFiles([]string{a, b})
	if err != nil {
		t.Fatalf("LoadEnvFiles: %v", err)
	}
	if m["X"] != "overridden" {
		t.Errorf("X: want overridden, got %q", m["X"])
	}
	if m["Y"] != "2" {
		t.Errorf("Y: want 2, got %q", m["Y"])
	}
}

func TestSubstituteEnv_ResolvesReferences(t *testing.T) {
	in := map[string]any{
		"exporters": map[string]any{
			"otlp": map[string]any{
				"endpoint": "${env:OTEL_URL}",
				"headers": map[string]any{
					"x-token":  "${ENV:TOKEN}",
					"x-static": "literal",
				},
				"tags": []any{"${env:REGION}", "static"},
			},
		},
	}
	env := map[string]string{
		"OTEL_URL": "https://collector.example.com",
		"TOKEN":    "abc123",
		"REGION":   "us-east-1",
	}
	config.SubstituteEnv(in, env)

	exporters := in["exporters"].(map[string]any)
	otlp := exporters["otlp"].(map[string]any)
	if got := otlp["endpoint"]; got != "https://collector.example.com" {
		t.Errorf("endpoint: got %v", got)
	}
	headers := otlp["headers"].(map[string]any)
	if got := headers["x-token"]; got != "abc123" {
		t.Errorf("x-token: got %v", got)
	}
	tags := otlp["tags"].([]any)
	if tags[0] != "us-east-1" {
		t.Errorf("tags[0]: got %v", tags[0])
	}
}

func TestSubstituteEnv_LeavesUnknownIntact(t *testing.T) {
	in := map[string]any{"endpoint": "${env:UNKNOWN}"}
	config.SubstituteEnv(in, map[string]string{"OTHER": "value"})
	if in["endpoint"] != "${env:UNKNOWN}" {
		t.Errorf("unknown var should be preserved, got %v", in["endpoint"])
	}
}

func TestSubstituteEnv_DefaultFallback(t *testing.T) {
	in := map[string]any{
		"endpoint": "${env:MISSING:-https://fallback.example.com}",
		"empty":    "${env:MISSING:-}",
	}
	config.SubstituteEnv(in, nil)
	if got := in["endpoint"]; got != "https://fallback.example.com" {
		t.Errorf("default fallback: got %v", got)
	}
	if got := in["empty"]; got != "" {
		t.Errorf("empty default: got %v", got)
	}
}

func TestSubstituteEnv_PartialString(t *testing.T) {
	in := map[string]any{"endpoint": "https://${env:HOST}:4318/v1/traces"}
	config.SubstituteEnv(in, map[string]string{"HOST": "collector"})
	if in["endpoint"] != "https://collector:4318/v1/traces" {
		t.Errorf("partial: got %v", in["endpoint"])
	}
}

func TestSubstituteEnv_NilEnvNoDefaultsIsNoop(t *testing.T) {
	in := map[string]any{"endpoint": "${env:HOST}"}
	config.SubstituteEnv(in, nil)
	if in["endpoint"] != "${env:HOST}" {
		t.Errorf("placeholder should remain when no env and no default: got %v", in["endpoint"])
	}
}

func TestLoadEnvFile_NotFound(t *testing.T) {
	if _, err := config.LoadEnvFile(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadEnvFile_LastWins(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "X=1\nX=2\n")
	m, err := config.LoadEnvFile(p)
	if err != nil {
		t.Fatalf("LoadEnvFile: %v", err)
	}
	if m["X"] != "2" {
		t.Errorf("last def should win, got %q", m["X"])
	}
}

func TestSubstituteEnv_NonStringScalarsUntouched(t *testing.T) {
	in := map[string]any{
		"port":    float64(4317),
		"enabled": true,
		"ratio":   0.5,
	}
	config.SubstituteEnv(in, map[string]string{"X": "y"})
	if in["port"] != float64(4317) || in["enabled"] != true || in["ratio"] != 0.5 {
		t.Errorf("scalars mutated: %v", in)
	}
}

// Ensure t.TempDir + os.WriteFile interplay is available via the existing
// writeFile helper in loader_test.go (same package).
var _ = os.WriteFile
