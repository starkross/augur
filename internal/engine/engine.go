package engine

import (
	"context"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
)

// Severity represents the severity level of a policy finding.
type Severity string

const (
	// SeverityDeny indicates a blocking violation that must be fixed.
	SeverityDeny Severity = "deny"
	// SeverityWarn indicates an advisory finding for best practices.
	SeverityWarn Severity = "warn"
)

// Finding represents a single policy violation found during evaluation.
type Finding struct {
	RuleID   string   `json:"rule_id"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	File     string   `json:"file"`
}

// Result holds all findings for a single evaluated file.
type Result struct {
	File     string    `json:"file"`
	Findings []Finding `json:"findings"`
}

// Engine evaluates OPA/Rego policies against OpenTelemetry Collector configs.
type Engine struct {
	pq rego.PreparedEvalQuery
}

// New compiles the Rego policies from the given filesystem and returns an Engine
// ready for evaluation. The policyDir specifies the root directory within the
// filesystem to walk for .rego files.
func New(policies fs.FS, policyDir string) (*Engine, error) {
	modules := make(map[string]string)

	err := fs.WalkDir(policies, policyDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".rego") || strings.HasSuffix(path, "_test.rego") {
			return nil
		}

		data, readErr := fs.ReadFile(policies, path)
		if readErr != nil {
			return fmt.Errorf("reading %q: %w", path, readErr)
		}

		modules[path] = string(data)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loading policies: %w", err)
	}

	compiler, compileErr := ast.CompileModules(modules)
	if compileErr != nil {
		return nil, fmt.Errorf("compiling policies: %w", compileErr)
	}

	pq, err := rego.New(
		rego.Query("data.main"),
		rego.Compiler(compiler),
	).PrepareForEval(context.Background())
	if err != nil {
		return nil, fmt.Errorf("preparing query: %w", err)
	}

	return &Engine{pq: pq}, nil
}

// Eval evaluates the compiled policies against the given input and returns
// the findings sorted by rule ID.
func (e *Engine) Eval(ctx context.Context, file string, input map[string]any) (*Result, error) {
	rs, err := e.pq.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, fmt.Errorf("evaluating policies: %w", err)
	}

	result := &Result{File: file}

	for _, r := range rs {
		for _, expr := range r.Expressions {
			obj, ok := expr.Value.(map[string]any)
			if !ok {
				continue
			}
			appendFindings(result, obj["deny"], SeverityDeny, file)
			appendFindings(result, obj["warn"], SeverityWarn, file)
		}
	}

	slices.SortFunc(result.Findings, func(a, b Finding) int {
		return strings.Compare(a.RuleID, b.RuleID)
	})

	return result, nil
}

func appendFindings(result *Result, val any, severity Severity, file string) {
	set, ok := val.([]any)
	if !ok {
		return
	}
	for _, item := range set {
		if s, ok := item.(string); ok {
			result.Findings = append(result.Findings, Finding{
				RuleID:   extractRuleID(s),
				Severity: severity,
				Message:  s,
				File:     file,
			})
		}
	}
}

func extractRuleID(msg string) string {
	if idx := strings.Index(msg, ":"); idx > 0 {
		candidate := strings.TrimSpace(msg[:idx])
		if strings.HasPrefix(candidate, "OTEL-") {
			return candidate
		}
	}
	return "UNKNOWN"
}
