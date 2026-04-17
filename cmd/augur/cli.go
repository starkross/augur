package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/starkross/augur"
	"github.com/starkross/augur/internal/output"
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
		RunE: func(_ *cobra.Command, args []string) error {
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
	skipRules []string
	outputFmt string
	policyDir string
	strict    bool
	quiet     bool
	noColor   bool
}

func run(files []string, opts runOpts) error {
	ctx := context.Background()

	linterOpts := []augur.Option{}
	if opts.policyDir != "" {
		linterOpts = append(linterOpts, augur.WithPolicyDir(opts.policyDir))
	}
	if len(opts.skipRules) > 0 {
		linterOpts = append(linterOpts, augur.WithSkipRules(opts.skipRules...))
	}
	if opts.quiet {
		linterOpts = append(linterOpts, augur.WithSeverities(augur.SeverityDeny))
	}

	linter, err := augur.New(linterOpts...)
	if err != nil {
		return fmt.Errorf("initializing linter: %w", err)
	}

	formatter, err := output.GetFormatter(opts.outputFmt, opts.noColor)
	if err != nil {
		return err
	}

	result, err := linter.LintFiles(ctx, files)
	if err != nil {
		return err
	}

	hasFailures := false
	for _, f := range result.Findings {
		if f.Severity == augur.SeverityDeny || (f.Severity == augur.SeverityWarn && opts.strict) {
			hasFailures = true
			break
		}
	}

	if err := formatter.Format(os.Stdout, []*augur.Result{result}); err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	if hasFailures {
		return errors.New("lint failures detected")
	}
	return nil
}

func parseSkip(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, id := range parts {
		id = strings.TrimSpace(id)
		if id != "" {
			out = append(out, id)
		}
	}
	return out
}
