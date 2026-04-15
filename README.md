---

<p align="center">
  <img src="docs/logo.png" alt="augur logo" width="300">
</p>

<p align="center">
  <strong>
    <a href="#install">Getting Started</a>
    &nbsp;&nbsp;&bull;&nbsp;&nbsp;
    <a href="CONTRIBUTING.md">Getting Involved</a>
    &nbsp;&nbsp;&bull;&nbsp;&nbsp;
    <a href="https://github.com/starkross/augur/issues">Getting In Touch</a>
  </strong>
</p>

<p align="center">
  <a href="https://github.com/starkross/augur/actions/workflows/ci.yml?query=branch%3Amain">
    <img alt="Build Status" src="https://img.shields.io/github/actions/workflow/status/starkross/augur/ci.yml?branch=main&style=for-the-badge">
  </a>
  <a href="https://goreportcard.com/report/github.com/starkross/augur">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/starkross/augur?style=for-the-badge">
  </a>
  <a href="https://github.com/starkross/augur/releases">
    <img alt="GitHub release (latest by date including pre-releases)" src="https://img.shields.io/github/v/release/starkross/augur?include_prereleases&style=for-the-badge">
  </a>
</p>

<p align="center">
  <strong>
    <a href="docs/RULES.md">Rules</a>
    &nbsp;&nbsp;&bull;&nbsp;&nbsp;
    <a href="#usage">Usage</a>
    &nbsp;&nbsp;&bull;&nbsp;&nbsp;
    <a href="#custom-policies">Custom Policies</a>
    &nbsp;&nbsp;&bull;&nbsp;&nbsp;
    <a href="https://pkg.go.dev/github.com/starkross/augur">Package</a>
  </strong>
</p>

---

# augur

A fast, opinionated linter for [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) configurations. Catches misconfigurations, security issues, and performance pitfalls before they hit production.

Built with [OPA/Rego](https://www.openpolicyagent.org/) — every rule is a plain `.rego` file you can read, override, or extend.

<p align="center">
  <img src="docs/demo.gif" alt="augur demo" width="900">
</p>

## Why

The OpenTelemetry Collector is flexible, but that flexibility makes it easy to ship configs that silently drop data, leak secrets, or OOM under load. `augur` encodes hard-won operational knowledge into automated checks:

- **No memory limiter?** You'll OOM in production.
- **Hardcoded API key?** It'll end up in version control.
- **Batch processor in the wrong position?** You're leaving performance on the table.

## Install

### Homebrew

```sh
brew install --cask starkross/tap/augur
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

See [docs/RULES.md](docs/RULES.md) for the full list of built-in rules.
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

### Multiple config files

When you pass more than one file, augur deep-merges them into a single effective config before linting — the same way the OpenTelemetry Collector combines multiple `--config` flags:

```sh
augur common.yaml client.yaml
```

Merge rules (matching the collector's confmap):

| Value type | Behavior |
|------------|----------|
| Map        | Deep-merged recursively. Later files add keys and override scalar leaves. |
| Scalar     | Later file wins. |
| Slice      | Replaced wholesale by the later file. |

To lint several standalone configs independently, invoke augur once per file.

### Examples

```sh
# Merge two files
augur gateway.yaml overrides.yaml

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

All built-in rules live in `rules/policy/` as standard Rego files. To add your own:

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

## Security

Every release is reproducibly built and signed:

- **CycloneDX SBOMs** for every archive and every container image
- **Cosign keyless signatures** (sigstore OIDC) on `checksums.txt` and on the published OCI image
- **SLSA Level 3 build provenance** for binary artifacts

Quick verification of a release:

```sh
cosign verify ghcr.io/starkross/augur:vX.Y.Z \
  --certificate-identity-regexp 'https://github.com/starkross/augur/\.github/workflows/release\.yml@refs/tags/v.*' \
  --certificate-oidc-issuer     'https://token.actions.githubusercontent.com'
```

See [SECURITY.md](SECURITY.md) for the full verification recipe (binaries, images, SLSA provenance, SBOMs) and the vulnerability disclosure process.
