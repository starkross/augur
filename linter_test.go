package augur_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/starkross/augur"
)

const goodYAML = `
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
processors:
  memory_limiter:
    check_interval: 5s
    limit_mib: 4000
  batch:
    send_batch_max_size: 16384
exporters:
  otlp/backend:
    endpoint: "${env:OTEL_EXPORTER_ENDPOINT}"
    tls:
      insecure: false
    retry_on_failure:
      enabled: true
    sending_queue:
      enabled: true
      storage: file_storage
extensions:
  health_check: {}
  file_storage:
    directory: /tmp/otel
service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/backend]
`

const badYAML = `
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
processors: {}
exporters:
  debug: {}
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug]
`

func TestNew_Default(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil linter")
	}
}

func TestLintYAML_GoodConfig(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintYAML(context.Background(), "good.yaml", []byte(goodYAML))
	if err != nil {
		t.Fatalf("LintYAML: %v", err)
	}

	for _, f := range result.Findings {
		if f.Severity == augur.SeverityDeny {
			t.Errorf("unexpected deny finding %s: %s", f.RuleID, f.Message)
		}
	}
}

func TestLintYAML_BadConfigHasDenials(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintYAML(context.Background(), "bad.yaml", []byte(badYAML))
	if err != nil {
		t.Fatalf("LintYAML: %v", err)
	}

	var denies int
	for _, f := range result.Findings {
		if f.Severity == augur.SeverityDeny {
			denies++
		}
	}
	if denies == 0 {
		t.Error("expected at least one deny finding for bad config")
	}
}

func TestLint_MapInput(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	input := map[string]any{
		"receivers":  map[string]any{"otlp": map[string]any{}},
		"processors": map[string]any{},
		"exporters":  map[string]any{"debug": map[string]any{}},
		"service": map[string]any{
			"pipelines": map[string]any{
				"traces": map[string]any{
					"receivers": []any{"otlp"},
					"exporters": []any{"debug"},
				},
			},
		},
	}

	result, err := l.Lint(context.Background(), "test.yaml", input)
	if err != nil {
		t.Fatalf("Lint: %v", err)
	}
	if result.File != "test.yaml" {
		t.Errorf("expected file=test.yaml, got %q", result.File)
	}
	if len(result.Findings) == 0 {
		t.Error("expected findings for minimal config")
	}
}

func TestLintFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(badYAML), 0o600); err != nil {
		t.Fatal(err)
	}

	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintFile(context.Background(), path)
	if err != nil {
		t.Fatalf("LintFile: %v", err)
	}
	if result.File != path {
		t.Errorf("expected file=%q, got %q", path, result.File)
	}
}

func TestLintFiles_DeepMerge(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.yaml")
	b := filepath.Join(dir, "b.yaml")
	if err := os.WriteFile(a, []byte(goodYAML), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("service:\n  telemetry:\n    logs:\n      level: debug\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintFiles(context.Background(), []string{a, b})
	if err != nil {
		t.Fatalf("LintFiles: %v", err)
	}

	var sawOTEL016 bool
	for _, f := range result.Findings {
		if f.RuleID == "OTEL-016" {
			sawOTEL016 = true
		}
	}
	if !sawOTEL016 {
		t.Error("expected OTEL-016 (debug log level) after merge")
	}
}

func TestLintFiles_EmptyPaths(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := l.LintFiles(context.Background(), nil); err == nil {
		t.Error("expected error for empty paths")
	}
}

func TestWithSkipRules(t *testing.T) {
	l, err := augur.New(augur.WithSkipRules("OTEL-001", "OTEL-003"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintYAML(context.Background(), "bad.yaml", []byte(badYAML))
	if err != nil {
		t.Fatalf("LintYAML: %v", err)
	}

	for _, f := range result.Findings {
		if f.RuleID == "OTEL-001" || f.RuleID == "OTEL-003" {
			t.Errorf("expected %s to be skipped", f.RuleID)
		}
	}
}

func TestWithSeverities_DenyOnly(t *testing.T) {
	l, err := augur.New(augur.WithSeverities(augur.SeverityDeny))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintYAML(context.Background(), "bad.yaml", []byte(badYAML))
	if err != nil {
		t.Fatalf("LintYAML: %v", err)
	}

	if len(result.Findings) == 0 {
		t.Fatal("expected at least one deny finding")
	}
	for _, f := range result.Findings {
		if f.Severity != augur.SeverityDeny {
			t.Errorf("expected only deny findings, got %s: %s", f.Severity, f.RuleID)
		}
	}
}

func TestWithPolicyFS_AddsCustomRule(t *testing.T) {
	custom := fstest.MapFS{
		"main/custom.rego": &fstest.MapFile{
			Data: []byte(`package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
	not input.extensions.pprof
	msg := "CUSTOM-001: pprof extension required by platform."
}
`),
		},
	}

	l, err := augur.New(augur.WithPolicyFS(custom, "."))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.LintYAML(context.Background(), "good.yaml", []byte(goodYAML))
	if err != nil {
		t.Fatalf("LintYAML: %v", err)
	}

	var sawCustom bool
	for _, f := range result.Findings {
		if f.RuleID == "CUSTOM-001" {
			sawCustom = true
		}
	}
	if !sawCustom {
		t.Error("expected custom rule to fire")
	}
}

func TestWithoutBuiltinRules_RequiresExtraSource(t *testing.T) {
	if _, err := augur.New(augur.WithoutBuiltinRules()); err == nil {
		t.Error("expected error when disabling built-ins without extra source")
	}
}

func TestWithoutBuiltinRules_OnlyCustom(t *testing.T) {
	custom := fstest.MapFS{
		"main/only.rego": &fstest.MapFile{
			Data: []byte(`package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
	true
	msg := "ONLY-001: fires always."
}
`),
		},
	}

	l, err := augur.New(
		augur.WithoutBuiltinRules(),
		augur.WithPolicyFS(custom, "."),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := l.Lint(context.Background(), "x.yaml", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("Lint: %v", err)
	}

	if len(result.Findings) == 0 {
		t.Fatal("expected custom finding")
	}
	for _, f := range result.Findings {
		if f.RuleID != "ONLY-001" {
			t.Errorf("expected only ONLY-001, got %s", f.RuleID)
		}
	}
}

func TestLintYAML_InvalidYAML(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := l.LintYAML(context.Background(), "bad", []byte(":\n  :\n    [invalid")); err == nil {
		t.Error("expected error for invalid YAML")
	}
}
