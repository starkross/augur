---
id: otel-056
title: "OTEL-056: filelog start_at:beginning without storage"
sidebar_label: OTEL-056
description: Re-reading from the beginning without a checkpoint causes duplicate ingestion on every restart.
---

# OTEL-056: `filelog` `start_at:beginning` without storage

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`start_at: beginning` tells `filelog` to read each new file from byte zero. Combined with no `storage` extension, every Collector restart re-reads every file — duplicating every log line already ingested since the file was created. In production this manifests as mysterious duplicate log spam after every rolling restart. Either use `start_at: end` or wire a persistent storage extension so `filelog` can checkpoint its read position.

This rule fires when a `filelog` receiver has `start_at: beginning` and no `storage` configured.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  filelog:
    include: [/var/log/app/*.log]
    start_at: beginning             # no storage checkpoint
```

:::

:::tip[Prefer]

```yaml
extensions:
  file_storage/filelog:
    directory: /var/lib/otelcol/filelog

receivers:
  filelog:
    include: [/var/log/app/*.log]
    start_at: beginning
    storage: file_storage/filelog

service:
  extensions: [file_storage/filelog]
```

:::

## When Not To Use It

Ephemeral containers where the file is guaranteed to exist only once and restart means the file is gone too — rare. In that case acknowledge the warning.

## Related Rules

- [OTEL-057](./otel-057) — `filelog` overly broad include pattern
- [OTEL-065](../reliability/otel-065) — `sending_queue` without persistent storage
- [OTEL-066](../extension/otel-066) — `sending_queue.storage` references undefined extension

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — filelog receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/filelogreceiver/README.md)
- [OpenTelemetry Collector Contrib — file_storage extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/storage/filestorage/README.md)

## Resources

- Rule source: [`policy/main/receiver.rego`](https://github.com/starkross/augur/blob/main/policy/main/receiver.rego)
