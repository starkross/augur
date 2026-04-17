package augur

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/starkross/augur/internal/config"
	"github.com/starkross/augur/internal/engine"
	"github.com/starkross/augur/rules"
)

// Severity is the severity level of a policy finding.
type Severity string

const (
	// SeverityDeny indicates a blocking violation that must be fixed.
	SeverityDeny Severity = "deny"
	// SeverityWarn indicates an advisory finding for best practices.
	SeverityWarn Severity = "warn"
)

// Finding represents a single policy violation.
type Finding struct {
	RuleID   string   `json:"rule_id"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	File     string   `json:"file"`
}

// Result holds the findings produced by a single lint evaluation.
type Result struct {
	File     string    `json:"file"`
	Findings []Finding `json:"findings"`
}

// Linter evaluates OpenTelemetry Collector configs against a compiled policy
// set. A Linter is safe for concurrent use.
type Linter struct {
	engine     *engine.Engine
	skipRules  map[string]struct{}
	severities map[Severity]struct{}
}

// New constructs a Linter. By default it loads the bundled OTEL-* rules; use
// [WithPolicyDir] or [WithPolicyFS] to add custom policies, and
// [WithoutBuiltinRules] to exclude the built-ins.
func New(opts ...Option) (*Linter, error) {
	var lo linterOptions
	for _, opt := range opts {
		opt(&lo)
	}

	sources := make([]engine.PolicySource, 0, 1+len(lo.extraSources))
	if !lo.disableBuiltins {
		sources = append(sources, engine.PolicySource{
			FS:  rules.Policies,
			Dir: rules.PolicyDir,
		})
	}
	sources = append(sources, lo.extraSources...)
	if len(sources) == 0 {
		return nil, errors.New("augur: no policies configured; WithoutBuiltinRules requires at least one WithPolicyDir or WithPolicyFS")
	}

	eng, err := engine.New(sources...)
	if err != nil {
		return nil, err
	}

	return &Linter{
		engine:     eng,
		skipRules:  lo.skipRules,
		severities: lo.severities,
	}, nil
}

// Lint evaluates a pre-parsed config map. label is used only for display
// (e.g., in Finding.File).
func (l *Linter) Lint(ctx context.Context, label string, input map[string]any) (*Result, error) {
	r, err := l.engine.Eval(ctx, label, input)
	if err != nil {
		return nil, err
	}
	return l.filter(fromEngineResult(r)), nil
}

// LintYAML parses raw YAML bytes and evaluates the resulting config.
func (l *Linter) LintYAML(ctx context.Context, label string, data []byte) (*Result, error) {
	m, err := config.ParseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}
	return l.Lint(ctx, label, m)
}

// LintFile reads and evaluates a single YAML file.
func (l *Linter) LintFile(ctx context.Context, path string) (*Result, error) {
	m, err := config.LoadYAML(path)
	if err != nil {
		return nil, err
	}
	return l.Lint(ctx, path, m)
}

// LintFiles deep-merges several YAML files and evaluates the result. Merge
// semantics match the OpenTelemetry Collector's --config behavior: maps
// merge recursively; scalars and slices are replaced by the later file.
func (l *Linter) LintFiles(ctx context.Context, paths []string) (*Result, error) {
	if len(paths) == 0 {
		return nil, errors.New("augur: no config files provided")
	}
	m, err := config.LoadMerged(paths)
	if err != nil {
		return nil, err
	}
	label := paths[0]
	if len(paths) > 1 {
		label = "merged: " + strings.Join(paths, ", ")
	}
	return l.Lint(ctx, label, m)
}

func fromEngineResult(r *engine.Result) *Result {
	findings := make([]Finding, len(r.Findings))
	for i, f := range r.Findings {
		findings[i] = Finding{
			RuleID:   f.RuleID,
			Severity: Severity(f.Severity),
			Message:  f.Message,
			File:     f.File,
		}
	}
	return &Result{File: r.File, Findings: findings}
}

func (l *Linter) filter(r *Result) *Result {
	if len(l.skipRules) == 0 && len(l.severities) == 0 {
		return r
	}
	out := &Result{File: r.File, Findings: make([]Finding, 0, len(r.Findings))}
	for _, f := range r.Findings {
		if _, skip := l.skipRules[f.RuleID]; skip {
			continue
		}
		if len(l.severities) > 0 {
			if _, keep := l.severities[f.Severity]; !keep {
				continue
			}
		}
		out.Findings = append(out.Findings, f)
	}
	return out
}
