package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/starkross/augur/internal/config"
	"github.com/starkross/augur/internal/engine"
	"github.com/starkross/augur/internal/output"
	"github.com/starkross/augur/rules"
)

func newRootCmd(version string) *cobra.Command {
	var (
		outputFmt string
		strict    bool
		quiet     bool
		skipRules string
		noColor   bool
		policyDir string
	)

	root := &cobra.Command{
		Use:   "augur [flags] <config.yaml> [config.yaml...]",
		Short: "Lint OpenTelemetry Collector configs for best practices",
		Long: "Lint OpenTelemetry Collector configs for best practices.\n\n" +
			"When multiple files are provided, they are deep-merged into a single " +
			"effective config before linting, matching the collector's own --config " +
			"behavior (maps merge recursively; slices and scalars are replaced by the later file).",
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, runOpts{
				outputFmt: outputFmt,
				strict:    strict,
				quiet:     quiet,
				skipRules: parseSkip(skipRules),
				noColor:   noColor,
				policyDir: policyDir,
			})
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	f := root.Flags()
	f.StringVarP(&outputFmt, "output", "o", "text", "Output format: text, json, github")
	f.BoolVarP(&strict, "strict", "s", false, "Treat warnings as errors")
	f.BoolVarP(&quiet, "quiet", "q", false, "Only show failures, suppress warnings")
	f.StringVarP(&skipRules, "skip", "k", "", "Comma-separated rule IDs to skip")
	f.BoolVar(&noColor, "no-color", false, "Disable colored output")
	f.StringVarP(&policyDir, "policy", "p", "", "Additional policy directory (merged with built-in rules)")

	return root
}

type runOpts struct {
	skipRules map[string]struct{}
	outputFmt string
	policyDir string
	strict    bool
	quiet     bool
	noColor   bool
}

func run(files []string, opts runOpts) error {
	ctx := context.Background()

	sources := []engine.PolicySource{
		{FS: rules.Policies, Dir: rules.PolicyDir},
	}
	if opts.policyDir != "" {
		sources = append(sources, engine.PolicySource{FS: os.DirFS(opts.policyDir), Dir: "."})
	}

	eng, err := engine.New(sources...)
	if err != nil {
		return fmt.Errorf("initializing policy engine: %w", err)
	}

	formatter, err := output.GetFormatter(opts.outputFmt, opts.noColor)
	if err != nil {
		return err
	}

	input, err := config.LoadMerged(files)
	if err != nil {
		return err
	}

	label := files[0]
	if len(files) > 1 {
		label = "merged: " + strings.Join(files, ", ")
	}

	result, evalErr := eng.Eval(ctx, label, input)
	if evalErr != nil {
		return fmt.Errorf("evaluating %q: %w", label, evalErr)
	}

	filtered := filterFindings(result, opts)
	hasFailures := false
	for _, finding := range filtered.Findings {
		if finding.Severity == engine.SeverityDeny || (finding.Severity == engine.SeverityWarn && opts.strict) {
			hasFailures = true
		}
	}

	if err := formatter.Format(os.Stdout, []*engine.Result{filtered}); err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	if hasFailures {
		return errors.New("lint failures detected")
	}
	return nil
}

func filterFindings(r *engine.Result, opts runOpts) *engine.Result {
	if len(opts.skipRules) == 0 && !opts.quiet {
		return r
	}
	filtered := &engine.Result{File: r.File, Findings: make([]engine.Finding, 0, len(r.Findings))}
	for _, finding := range r.Findings {
		if _, skip := opts.skipRules[finding.RuleID]; skip {
			continue
		}
		if opts.quiet && finding.Severity == engine.SeverityWarn {
			continue
		}
		filtered.Findings = append(filtered.Findings, finding)
	}
	return filtered
}

func parseSkip(s string) map[string]struct{} {
	m := make(map[string]struct{})
	if s == "" {
		return m
	}
	for _, id := range strings.Split(s, ",") {
		m[strings.TrimSpace(id)] = struct{}{}
	}
	return m
}
