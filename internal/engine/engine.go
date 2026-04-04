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

// PolicySource represents a filesystem and root directory to load policies from.
type PolicySource struct {
	FS  fs.FS
	Dir string
}

// New compiles the Rego policies from one or more sources and returns an Engine
// ready for evaluation. Modules from all sources are merged; later sources can
// extend (but not override) earlier ones.
func New(sources ...PolicySource) (*Engine, error) {
	modules := make(map[string]string)

	for i, src := range sources {
		prefix := fmt.Sprintf("src%d/", i)
		err := fs.WalkDir(src.FS, src.Dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".rego") || strings.HasSuffix(path, "_test.rego") {
				return nil
			}

			data, readErr := fs.ReadFile(src.FS, path)
			if readErr != nil {
				return fmt.Errorf("reading %q: %w", path, readErr)
			}

			modules[prefix+path] = string(data)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("loading policies: %w", err)
		}
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
