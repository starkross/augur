package engine_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/starkross/augur/internal/engine"
	"github.com/starkross/augur/internal/rules"
)

func TestNew_Embedded(t *testing.T) {
	_, err := engine.New(engine.PolicySource{FS: rules.Policies, Dir: rules.PolicyDir})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
}

func TestNew_CustomFS(t *testing.T) {
	policyFS := fstest.MapFS{
		"main/main.rego": &fstest.MapFile{
			Data: []byte(`package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
	not input.valid
	msg := "TEST: input is not valid"
}
`),
		},
	}

	eng, err := engine.New(engine.PolicySource{FS: policyFS, Dir: "."})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	result, err := eng.Eval(context.Background(), "test.yaml", map[string]any{"valid": false})
	if err != nil {
		t.Fatalf("Eval() returned error: %v", err)
	}

	if len(result.Findings) == 0 {
		t.Error("expected findings for invalid input")
	}
}

func TestNew_InvalidRego(t *testing.T) {
	policyFS := fstest.MapFS{
		"bad.rego": &fstest.MapFile{Data: []byte("not valid rego {{{{")},
	}

	_, err := engine.New(engine.PolicySource{FS: policyFS, Dir: "."})
	if err == nil {
		t.Fatal("expected error for invalid Rego")
	}
}

func TestEval_DenyFindings(t *testing.T) {
	eng, err := engine.New(engine.PolicySource{FS: rules.Policies, Dir: rules.PolicyDir})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
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

	result, err := eng.Eval(context.Background(), "test.yaml", input)
	if err != nil {
		t.Fatalf("Eval() returned error: %v", err)
	}

	var denyCount int
	for _, f := range result.Findings {
		if f.Severity == engine.SeverityDeny {
			denyCount++
		}
	}
	if denyCount == 0 {
		t.Error("expected at least one deny finding for minimal config")
	}
}

func TestEval_NoFindings(t *testing.T) {
	eng, err := engine.New(engine.PolicySource{FS: rules.Policies, Dir: rules.PolicyDir})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	input := map[string]any{
		"receivers": map[string]any{
			"otlp": map[string]any{
				"protocols": map[string]any{
					"grpc": map[string]any{"endpoint": "localhost:4317"},
				},
			},
		},
		"processors": map[string]any{
			"batch":          map[string]any{"send_batch_max_size": 16384},
			"memory_limiter": map[string]any{"check_interval": "5s", "limit_mib": 4000},
		},
		"exporters": map[string]any{
			"otlp/backend": map[string]any{
				"endpoint":         "${env:OTEL_EXPORTER_ENDPOINT}",
				"tls":              map[string]any{"insecure": false},
				"retry_on_failure": map[string]any{"enabled": true},
				"sending_queue":    map[string]any{"enabled": true, "storage": "file_storage"},
			},
		},
		"extensions": map[string]any{
			"health_check": map[string]any{},
			"file_storage": map[string]any{"directory": "/tmp/otel"},
		},
		"service": map[string]any{
			"extensions": []any{"health_check"},
			"pipelines": map[string]any{
				"traces": map[string]any{
					"receivers":  []any{"otlp"},
					"processors": []any{"memory_limiter", "batch"},
					"exporters":  []any{"otlp/backend"},
				},
			},
		},
	}

	result, err := eng.Eval(context.Background(), "good.yaml", input)
	if err != nil {
		t.Fatalf("Eval() returned error: %v", err)
	}

	if len(result.Findings) != 0 {
		t.Errorf("expected no findings, got %d:", len(result.Findings))
		for _, f := range result.Findings {
			t.Logf("  %s: %s", f.RuleID, f.Message)
		}
	}
}

func TestEval_FindingsSortedByRuleID(t *testing.T) {
	eng, err := engine.New(engine.PolicySource{FS: rules.Policies, Dir: rules.PolicyDir})
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
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

	result, err := eng.Eval(context.Background(), "test.yaml", input)
	if err != nil {
		t.Fatalf("Eval() returned error: %v", err)
	}

	for i := 1; i < len(result.Findings); i++ {
		if result.Findings[i].RuleID < result.Findings[i-1].RuleID {
			t.Errorf("findings not sorted: %s comes after %s",
				result.Findings[i].RuleID, result.Findings[i-1].RuleID)
		}
	}
}
