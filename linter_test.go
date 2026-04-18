package augur_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/starkross/augur"
)

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return data
}

const customPolicy = `package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
	not input.extensions.pprof
	msg := "CUSTOM-001: pprof extension required by platform."
}
`

const alwaysFiringPolicy = `package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
	true
	msg := "ONLY-001: fires always."
}
`

func TestLintYAML(t *testing.T) {
	customFS := fstest.MapFS{"main/custom.rego": &fstest.MapFile{Data: []byte(customPolicy)}}

	tests := []struct {
		name   string
		opts   []augur.Option
		file   string
		assert func(t *testing.T, r *augur.Result)
	}{
		{
			name: "good config has no deny findings",
			file: "good.yaml",
			assert: func(t *testing.T, r *augur.Result) {
				for _, f := range r.Findings {
					if f.Severity == augur.SeverityDeny {
						t.Errorf("unexpected deny %s: %s", f.RuleID, f.Message)
					}
				}
			},
		},
		{
			name: "bad config yields deny findings",
			file: "bad.yaml",
			assert: func(t *testing.T, r *augur.Result) {
				if !slices.ContainsFunc(r.Findings, func(f augur.Finding) bool {
					return f.Severity == augur.SeverityDeny
				}) {
					t.Error("expected at least one deny finding")
				}
			},
		},
		{
			name: "WithSkipRules drops matching rule IDs",
			opts: []augur.Option{augur.WithSkipRules("OTEL-001", "OTEL-003")},
			file: "bad.yaml",
			assert: func(t *testing.T, r *augur.Result) {
				for _, f := range r.Findings {
					if f.RuleID == "OTEL-001" || f.RuleID == "OTEL-003" {
						t.Errorf("expected %s to be skipped", f.RuleID)
					}
				}
			},
		},
		{
			name: "WithSeverities deny-only filters warnings",
			opts: []augur.Option{augur.WithSeverities(augur.SeverityDeny)},
			file: "bad.yaml",
			assert: func(t *testing.T, r *augur.Result) {
				if len(r.Findings) == 0 {
					t.Fatal("expected at least one deny finding")
				}
				for _, f := range r.Findings {
					if f.Severity != augur.SeverityDeny {
						t.Errorf("unexpected severity %s for %s", f.Severity, f.RuleID)
					}
				}
			},
		},
		{
			name: "WithPolicyFS adds custom rule",
			opts: []augur.Option{augur.WithPolicyFS(customFS, ".")},
			file: "good.yaml",
			assert: func(t *testing.T, r *augur.Result) {
				if !slices.ContainsFunc(r.Findings, func(f augur.Finding) bool {
					return f.RuleID == "CUSTOM-001"
				}) {
					t.Error("expected custom rule to fire")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l, err := augur.New(tc.opts...)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			r, err := l.LintYAML(context.Background(), tc.file, readTestdata(t, tc.file))
			if err != nil {
				t.Fatalf("LintYAML: %v", err)
			}
			tc.assert(t, r)
		})
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

	r, err := l.Lint(context.Background(), "test.yaml", input)
	if err != nil {
		t.Fatalf("Lint: %v", err)
	}
	if r.File != "test.yaml" {
		t.Errorf("expected file=test.yaml, got %q", r.File)
	}
	if len(r.Findings) == 0 {
		t.Error("expected findings for minimal config")
	}
}

func TestLintFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, readTestdata(t, "bad.yaml"), 0o600); err != nil {
		t.Fatal(err)
	}

	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	r, err := l.LintFile(context.Background(), path)
	if err != nil {
		t.Fatalf("LintFile: %v", err)
	}
	if r.File != path {
		t.Errorf("expected file=%q, got %q", path, r.File)
	}
}

func TestLintFiles_DeepMerge(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.yaml")
	b := filepath.Join(dir, "b.yaml")
	if err := os.WriteFile(a, readTestdata(t, "good.yaml"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, readTestdata(t, "debug_level.yaml"), 0o600); err != nil {
		t.Fatal(err)
	}

	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	r, err := l.LintFiles(context.Background(), []string{a, b})
	if err != nil {
		t.Fatalf("LintFiles: %v", err)
	}

	if !slices.ContainsFunc(r.Findings, func(f augur.Finding) bool {
		return f.RuleID == "OTEL-016"
	}) {
		t.Error("expected OTEL-016 (debug log level) after merge")
	}
}

func TestNew_Errors(t *testing.T) {
	tests := []struct {
		name   string
		opts   []augur.Option
		wantIs error
	}{
		{
			name:   "WithoutBuiltinRules requires extra source",
			opts:   []augur.Option{augur.WithoutBuiltinRules()},
			wantIs: augur.ErrNoPolicies,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := augur.New(tc.opts...)
			if err == nil {
				t.Fatal("expected error")
			}
			if tc.wantIs != nil && !errors.Is(err, tc.wantIs) {
				t.Errorf("errors.Is(%v, %v) = false", err, tc.wantIs)
			}
		})
	}
}

func TestLint_Errors(t *testing.T) {
	l, err := augur.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx := context.Background()

	t.Run("LintFiles empty paths", func(t *testing.T) {
		if _, err := l.LintFiles(ctx, nil); err == nil {
			t.Error("expected error for empty paths")
		}
	})

	t.Run("LintYAML invalid YAML", func(t *testing.T) {
		if _, err := l.LintYAML(ctx, "bad", readTestdata(t, "invalid.yaml")); err == nil {
			t.Error("expected error for invalid YAML")
		}
	})
}

func TestWithoutBuiltinRules_OnlyCustom(t *testing.T) {
	custom := fstest.MapFS{"main/only.rego": &fstest.MapFile{Data: []byte(alwaysFiringPolicy)}}

	l, err := augur.New(
		augur.WithoutBuiltinRules(),
		augur.WithPolicyFS(custom, "."),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	r, err := l.Lint(context.Background(), "x.yaml", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("Lint: %v", err)
	}

	if len(r.Findings) == 0 {
		t.Fatal("expected custom finding")
	}
	for _, f := range r.Findings {
		if f.RuleID != "ONLY-001" {
			t.Errorf("expected only ONLY-001, got %s", f.RuleID)
		}
	}
}
