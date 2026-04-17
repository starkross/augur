// Package augur lints OpenTelemetry Collector YAML configurations against a
// set of OPA/Rego policies.
//
// It can be used as a command-line tool (see cmd/augur) or embedded as a Go
// library.
//
// # Library usage
//
// Construct a [Linter] with [New] and call one of the Lint methods:
//
//	linter, err := augur.New()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	result, err := linter.LintFile(context.Background(), "otel-collector.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, f := range result.Findings {
//	    fmt.Printf("%s [%s] %s\n", f.RuleID, f.Severity, f.Message)
//	}
//
// By default [New] loads the bundled OTEL-* rules. Supply additional policies
// with [WithPolicyDir] or [WithPolicyFS], or disable the built-ins with
// [WithoutBuiltinRules].
//
// Linters are safe for concurrent use after construction. A single compiled
// policy engine is reused across all calls.
package augur
