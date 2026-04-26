package augur

import (
	"io/fs"
	"os"

	"github.com/starkross/augur/internal/engine"
)

// Option configures a [Linter].
type Option func(*linterOptions)

type linterOptions struct {
	skipRules       map[string]struct{}
	severities      map[Severity]struct{}
	extraSources    []engine.PolicySource
	disableBuiltins bool
}

// WithPolicyDir adds a directory of custom .rego policies. The directory is
// walked recursively and merged with the built-in rules unless
// [WithoutBuiltinRules] is also supplied.
func WithPolicyDir(path string) Option {
	return func(o *linterOptions) {
		o.extraSources = append(o.extraSources, engine.PolicySource{
			FS:  os.DirFS(path),
			Dir: ".",
		})
	}
}

// WithPolicyFS adds a custom policy source backed by an [io/fs.FS]. root is
// the subdirectory within fsys that contains the .rego files — pass "." for
// the filesystem root. Useful for embedding additional policies into your
// own application binary via //go:embed.
func WithPolicyFS(fsys fs.FS, root string) Option {
	return func(o *linterOptions) {
		if root == "" {
			root = "."
		}
		o.extraSources = append(o.extraSources, engine.PolicySource{
			FS:  fsys,
			Dir: root,
		})
	}
}

// WithoutBuiltinRules disables the bundled OTEL-* rule set. When used, at
// least one of [WithPolicyDir] or [WithPolicyFS] must also be supplied.
func WithoutBuiltinRules() Option {
	return func(o *linterOptions) {
		o.disableBuiltins = true
	}
}

// WithSkipRules drops findings whose rule ID matches any of the given IDs.
// Empty strings are ignored. Multiple calls accumulate: the final skip set is
// the union of all IDs passed across every WithSkipRules call.
func WithSkipRules(ids ...string) Option {
	return func(o *linterOptions) {
		if o.skipRules == nil {
			o.skipRules = make(map[string]struct{}, len(ids))
		}
		for _, id := range ids {
			if id != "" {
				o.skipRules[id] = struct{}{}
			}
		}
	}
}

// WithSeverities restricts findings to the given severities. If not called,
// all severities are returned. Multiple calls accumulate: the final allow-set
// is the union of severities passed across every WithSeverities call.
func WithSeverities(severities ...Severity) Option {
	return func(o *linterOptions) {
		if o.severities == nil {
			o.severities = make(map[Severity]struct{}, len(severities))
		}
		for _, s := range severities {
			o.severities[s] = struct{}{}
		}
	}
}
