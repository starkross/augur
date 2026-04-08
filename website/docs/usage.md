---
id: usage
title: Usage
sidebar_position: 4
---

# Usage

```text
augur [flags] <config.yaml> [config.yaml...]
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output format: `text`, `json`, `github` | `text` |
| `-s, --strict` | Treat warnings as errors | `false` |
| `-q, --quiet` | Suppress warnings, show only failures | `false` |
| `-k, --skip` | Comma-separated rule IDs to skip | |
| `--no-color` | Disable colored output | `false` |
| `-p, --policy` | Additional policy directory (merged with built-in rules) | |

## Examples

Lint multiple files:

```sh
augur gateway.yaml agent.yaml
```

Strict mode — warnings become errors:

```sh
augur --strict config.yaml
```

JSON output for programmatic consumption:

```sh
augur -o json config.yaml
```

GitHub Actions annotation output:

```sh
augur -o github config.yaml
```

Skip specific rules:

```sh
augur --skip OTEL-015,OTEL-016 config.yaml
```

Merge in custom policies alongside built-in rules:

```sh
augur --policy ./my-policies config.yaml
```
