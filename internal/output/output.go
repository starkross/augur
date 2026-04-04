package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/starkross/augur/internal/engine"
)

// Formatter formats lint results and writes them to the given writer.
type Formatter interface {
	Format(w io.Writer, results []*engine.Result) error
}

// TextFormatter writes human-readable colored output.
type TextFormatter struct{ NoColor bool }

// Format writes results as colored text to w.
func (f *TextFormatter) Format(w io.Writer, results []*engine.Result) error {
	var werr error
	printf := func(format string, a ...any) {
		if werr == nil {
			_, werr = fmt.Fprintf(w, format, a...)
		}
	}

	denies, warns := 0, 0
	for _, r := range results {
		if len(r.Findings) == 0 {
			continue
		}
		printf("%s\n", f.bold(r.File))
		for _, finding := range r.Findings {
			switch finding.Severity {
			case engine.SeverityDeny:
				printf("  %s %s\n", f.red("FAIL"), finding.Message)
				denies++
			case engine.SeverityWarn:
				printf("  %s %s\n", f.yellow("WARN"), finding.Message)
				warns++
			}
		}
		printf("\n")
	}
	if denies+warns == 0 {
		printf("%s\n", f.green("✓ All checks passed"))
	} else if denies == 0 {
		printf("%s\n", f.yellow(fmt.Sprintf("⚠ %d warning(s), 0 failure(s)", warns)))
	} else {
		printf("%s\n", f.red(fmt.Sprintf("✗ %d failure(s), %d warning(s)", denies, warns)))
	}
	return werr
}

func (f *TextFormatter) red(s string) string {
	if f.NoColor {
		return s
	}
	return "\033[0;31m" + s + "\033[0m"
}

func (f *TextFormatter) yellow(s string) string {
	if f.NoColor {
		return s
	}
	return "\033[0;33m" + s + "\033[0m"
}

func (f *TextFormatter) green(s string) string {
	if f.NoColor {
		return s
	}
	return "\033[0;32m" + s + "\033[0m"
}

func (f *TextFormatter) bold(s string) string {
	if f.NoColor {
		return s
	}
	return "\033[1m" + s + "\033[0m"
}

// JSONFormatter writes results as structured JSON.
type JSONFormatter struct{}

type jsonOutput struct {
	Files   []*engine.Result `json:"files"`
	Summary jsonSummary      `json:"summary"`
}

type jsonSummary struct {
	TotalFiles int `json:"total_files"`
	Failures   int `json:"failures"`
	Warnings   int `json:"warnings"`
	Passed     int `json:"passed"`
}

// Format writes results as indented JSON to w.
func (f *JSONFormatter) Format(w io.Writer, results []*engine.Result) error {
	out := jsonOutput{Files: results}
	for _, r := range results {
		hasIssue := false
		for _, finding := range r.Findings {
			switch finding.Severity {
			case engine.SeverityDeny:
				out.Summary.Failures++
				hasIssue = true
			case engine.SeverityWarn:
				out.Summary.Warnings++
				hasIssue = true
			}
		}
		out.Summary.TotalFiles++
		if !hasIssue {
			out.Summary.Passed++
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// GitHubFormatter writes results as GitHub Actions workflow commands.
type GitHubFormatter struct{}

// Format writes results as ::error and ::warning annotations to w.
func (f *GitHubFormatter) Format(w io.Writer, results []*engine.Result) error {
	for _, r := range results {
		for _, finding := range r.Findings {
			level := "warning"
			if finding.Severity == engine.SeverityDeny {
				level = "error"
			}
			if _, err := fmt.Fprintf(w, "::%s file=%s,title=%s::%s\n",
				level, r.File, finding.RuleID, escapeGH(finding.Message)); err != nil {
				return err
			}
		}
	}
	return nil
}

func escapeGH(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\n", "%0A")
	s = strings.ReplaceAll(s, "\r", "%0D")
	return s
}

// GetFormatter returns a Formatter for the given format name.
// Supported formats: "text", "json", "github".
func GetFormatter(name string, noColor bool) (Formatter, error) {
	switch name {
	case "text":
		return &TextFormatter{NoColor: noColor}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "github":
		return &GitHubFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown output format: %q", name)
	}
}
