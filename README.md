# augur

A fast, opinionated linter for [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) configurations. Catches misconfigurations, security issues, and performance pitfalls before they hit production.

Built with [OPA/Rego](https://www.openpolicyagent.org/) — every rule is a plain `.rego` file you can read, override, or extend.

## Why

The OpenTelemetry Collector is flexible, but that flexibility makes it easy to ship configs that silently drop data, leak secrets, or OOM under load. `augur` encodes hard-won operational knowledge into automated checks:

- **No memory limiter?** You'll OOM in production.
- **Hardcoded API key?** It'll end up in version control.
- **Batch processor in the wrong position?** You're leaving performance on the table.

## Install

### Homebrew

```sh
brew install starkross/tap/augur
```

### Go

```sh
go install github.com/starkross/augur/cmd/augur@latest
```

### Docker

```sh
docker run --rm -v "$(pwd):/work" ghcr.io/starkross/augur:latest config.yaml
```

### Binary releases

Download from [GitHub Releases](https://github.com/starkross/augur/releases) — available for Linux, macOS, and Windows (amd64/arm64).

## Quick start

```sh
augur otel-collector-config.yaml
```

```
otel-collector-config.yaml
  FAIL OTEL-001: memory_limiter processor is not configured. Required to prevent OOM in production.
  FAIL OTEL-003: batch processor is not configured. Required for efficient data export.
  WARN OTEL-011: health_check extension is not configured. Recommended for k8s liveness/readiness probes.

✗ 2 failure(s), 1 warning(s)
```

Exit code `1` on any failure. Warnings are informational by default.

## Rules

### Deny (blocking)

| ID | Description |
|----|-------------|
| OTEL-001 | `memory_limiter` processor must be configured |
| OTEL-002 | `memory_limiter` must be included in every pipeline |
| OTEL-003 | `batch` processor must be configured |
| OTEL-004 | No hardcoded secrets in exporters |
| OTEL-005 | No hardcoded secrets in receivers |
| OTEL-006 | `service.pipelines` must be defined |
| OTEL-007 | Every pipeline must have receivers and exporters |

### Warn (advisory)

| ID | Description |
|----|-------------|
| OTEL-010 | Receivers should not bind to `0.0.0.0` |
| OTEL-011 | `health_check` extension recommended |
| OTEL-012 | `health_check` configured but not listed in `service.extensions` |
| OTEL-013 | `batch` processor should be last in pipeline |
| OTEL-014 | `memory_limiter` should be first processor in pipeline |
| OTEL-015 | `debug`/`logging` exporter detected |
| OTEL-016 | Telemetry log level set to `debug` |
| OTEL-017 | Exporter missing `retry_on_failure`/`sending_queue` |
| OTEL-018 | OTLP exporter without TLS on non-local endpoint |
| OTEL-020 | Unused receiver |
| OTEL-021 | Unused exporter |
| OTEL-022 | Unused processor |

List all rules from the CLI:

```sh
augur list-rules
```

## Usage

```
augur [flags] <config.yaml> [config.yaml...]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output format: `text`, `json`, `github` | `text` |
| `-s, --strict` | Treat warnings as errors | `false` |
| `-q, --quiet` | Suppress warnings, show only failures | `false` |
| `-k, --skip` | Comma-separated rule IDs to skip | |
| `--no-color` | Disable colored output | `false` |
| `-p, --policy` | Additional policy directory (merged with built-in rules) | |

### Examples

```sh
# Lint multiple files
augur gateway.yaml agent.yaml

# Strict mode — warnings become errors
augur --strict config.yaml

# JSON output for programmatic consumption
augur -o json config.yaml

# Skip specific rules
augur --skip OTEL-015,OTEL-016 config.yaml

# Use custom policies
augur --policy ./my-policies config.yaml
```

## Custom policies

All built-in rules live in `policy/` as standard Rego files. To add your own:

1. Create a directory with your custom rules:

```rego
# my-policies/main/custom.rego
package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
    not input.processors.filter
    msg := "CUSTOM-001: filter processor is required by our platform team."
}
```

2. Run with `--policy`:

```sh
augur --policy ./my-policies config.yaml
```

Custom policies are **merged** with the built-in rules — your rules run alongside all default checks.
