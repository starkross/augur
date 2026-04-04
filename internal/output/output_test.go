package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/starkross/augur/internal/engine"
	"github.com/starkross/augur/internal/output"
)

func sampleResults() []*engine.Result {
	return []*engine.Result{
		{
			File: "bad.yaml",
			Findings: []engine.Finding{
				{RuleID: "OTEL-001", Severity: engine.SeverityDeny, Message: "OTEL-001: missing memory_limiter", File: "bad.yaml"},
				{RuleID: "OTEL-010", Severity: engine.SeverityWarn, Message: "OTEL-010: binds to 0.0.0.0", File: "bad.yaml"},
			},
		},
	}
}

func emptyResults() []*engine.Result {
	return []*engine.Result{
		{File: "good.yaml"},
	}
}

func TestTextFormatter_WithFindings(t *testing.T) {
	var buf bytes.Buffer
	f := &output.TextFormatter{NoColor: true}
	if err := f.Format(&buf, sampleResults()); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "FAIL") {
		t.Error("expected FAIL in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(out, "1 failure(s)") {
		t.Error("expected failure summary")
	}
}

func TestTextFormatter_AllPassed(t *testing.T) {
	var buf bytes.Buffer
	f := &output.TextFormatter{NoColor: true}
	if err := f.Format(&buf, emptyResults()); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	if !strings.Contains(buf.String(), "All checks passed") {
		t.Error("expected 'All checks passed' in output")
	}
}

func TestTextFormatter_WarningsOnly(t *testing.T) {
	results := []*engine.Result{
		{
			File: "warn.yaml",
			Findings: []engine.Finding{
				{RuleID: "OTEL-010", Severity: engine.SeverityWarn, Message: "OTEL-010: warning", File: "warn.yaml"},
			},
		},
	}

	var buf bytes.Buffer
	f := &output.TextFormatter{NoColor: true}
	if err := f.Format(&buf, results); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "1 warning(s), 0 failure(s)") {
		t.Errorf("expected warning summary, got: %s", out)
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &output.JSONFormatter{}
	if err := f.Format(&buf, sampleResults()); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	summary, ok := parsed["summary"].(map[string]any)
	if !ok {
		t.Fatal("missing summary in JSON output")
	}
	if failures := summary["failures"]; failures != float64(1) {
		t.Errorf("expected 1 failure, got %v", failures)
	}
	if warnings := summary["warnings"]; warnings != float64(1) {
		t.Errorf("expected 1 warning, got %v", warnings)
	}
}

func TestGitHubFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := &output.GitHubFormatter{}
	if err := f.Format(&buf, sampleResults()); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::error") {
		t.Error("expected ::error annotation")
	}
	if !strings.Contains(out, "::warning") {
		t.Error("expected ::warning annotation")
	}
}

func TestGetFormatter_Valid(t *testing.T) {
	for _, name := range []string{"text", "json", "github"} {
		f, err := output.GetFormatter(name, false)
		if err != nil {
			t.Errorf("GetFormatter(%q) error: %v", name, err)
		}
		if f == nil {
			t.Errorf("GetFormatter(%q) returned nil", name)
		}
	}
}

func TestGetFormatter_Unknown(t *testing.T) {
	_, err := output.GetFormatter("xml", false)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
