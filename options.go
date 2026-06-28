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
	env             map[string]string
	envFiles        []string
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

// WithEnv supplies variable values for ${env:VAR} / ${ENV:VAR} substitution
// inside config string values. Substitution happens after YAML parsing and
// before policy evaluation, so rules see the resolved values. Unknown
// references are left untouched and remain detectable by the bundled
// `is_env_var` helper. Multiple calls merge — later calls override earlier
// keys.
func WithEnv(env map[string]string) Option {
	return func(o *linterOptions) {
		if o.env == nil {
			o.env = make(map[string]string, len(env))
		}
		for k, v := range env {
			o.env[k] = v
		}
	}
}

// WithEnvFile registers a .env-style file (KEY=VALUE per line; '#' comments)
// to be loaded at [New] time. The merged map drives ${env:VAR} substitution
// in the same way as [WithEnv]. Files are loaded in the order they are
// declared; later files override earlier ones, and explicit [WithEnv] entries
// override values from any file.
func WithEnvFile(paths ...string) Option {
	return func(o *linterOptions) {
		o.envFiles = append(o.envFiles, paths...)
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
