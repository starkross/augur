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
type Severity = engine.Severity

// Finding represents a single policy violation.
type Finding = engine.Finding

// Result holds the findings produced by a single lint evaluation.
type Result = engine.Result

const (
	// SeverityDeny indicates a blocking violation that must be fixed.
	SeverityDeny = engine.SeverityDeny
	// SeverityWarn indicates an advisory finding for best practices.
	SeverityWarn = engine.SeverityWarn
)

// ErrNoPolicies is returned by [New] when the built-in rules are disabled
// via [WithoutBuiltinRules] and no custom policy source is supplied.
var ErrNoPolicies = errors.New("no policies configured")

// Linter evaluates OpenTelemetry Collector configs against a compiled policy
// set. A Linter is safe for concurrent use.
type Linter struct {
	engine     *engine.Engine
	skipRules  map[string]struct{}
	severities map[Severity]struct{}
	env        map[string]string
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
		return nil, fmt.Errorf("%w: WithoutBuiltinRules requires at least one WithPolicyDir or WithPolicyFS", ErrNoPolicies)
	}

	eng, err := engine.New(sources...)
	if err != nil {
		return nil, fmt.Errorf("augur: compiling policies: %w", err)
	}

	env, err := buildEnv(&lo)
	if err != nil {
		return nil, err
	}

	return &Linter{
		engine:     eng,
		skipRules:  lo.skipRules,
		severities: lo.severities,
		env:        env,
	}, nil
}

func buildEnv(lo *linterOptions) (map[string]string, error) {
	if len(lo.envFiles) == 0 && len(lo.env) == 0 {
		return nil, nil
	}
	merged, err := config.LoadEnvFiles(lo.envFiles)
	if err != nil {
		return nil, fmt.Errorf("augur: loading env file: %w", err)
	}
	if merged == nil {
		merged = map[string]string{}
	}
	for k, v := range lo.env {
		merged[k] = v
	}
	return merged, nil
}

// Lint evaluates a pre-parsed config map. label identifies the source in
// the returned [Result.File] and each [Finding.File]; pass a filename, URL,
// or any other human-readable identifier.
func (l *Linter) Lint(ctx context.Context, label string, input map[string]any) (*Result, error) {
	config.SubstituteEnv(input, l.env)
	r, err := l.engine.Eval(ctx, label, input)
	if err != nil {
		return nil, err
	}
	return l.filter(r), nil
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

// LintFiles deep-merges the given YAML files and evaluates the result,
// matching the OpenTelemetry Collector's --config behavior.
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
