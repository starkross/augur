---
id: otel-057
title: "OTEL-057: filelog overly broad include pattern"
sidebar_label: OTEL-057
description: A /** glob ingests everything under the directory, including rotated files and binary blobs.
---

# OTEL-057: `filelog` overly broad include pattern

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Patterns like `/var/log/**/*` or `/var/log/**` tell `filelog` to watch every file under a tree — rotated archives, compressed logs, other services' directories, sometimes binary files the filesystem watcher has to skip. Scope the pattern to the specific files you actually want (`/var/log/myapp/*.log`) to avoid wasted CPU, wasted memory, and accidental ingestion of unrelated data.

This rule fires when any entry in `filelog.include` ends with `/*`, `/**`, or `/**/*`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  filelog:
    include:
      - /var/log/**                    # far too broad
    start_at: end
```

:::

:::tip[Prefer]

```yaml
receivers:
  filelog:
    include:
      - /var/log/myapp/app.log
      - /var/log/myapp/error.log
    exclude:
      - /var/log/myapp/*.gz
    start_at: end
```

:::

## When Not To Use It

You genuinely want everything in a directory and the files really are homogeneous. Even then, list the extensions explicitly.

## Related Rules

- [OTEL-056](./otel-056) — `filelog` `start_at:beginning` without storage

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — filelog receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/filelogreceiver/README.md)

## Resources

- Rule source: [`policy/main/receiver.rego`](https://github.com/starkross/augur/blob/main/policy/main/receiver.rego)
