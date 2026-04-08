---
id: otel-066
title: "OTEL-066: sending_queue.storage references undefined extension"
sidebar_label: OTEL-066
description: A queue storage reference has to match a real storage extension.
---

# OTEL-066: `sending_queue.storage` references undefined extension

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

When you enable persistent queue storage on an exporter with `sending_queue.storage: file_storage/queue`, the value has to match an extension name that is actually declared in the `extensions:` block. A typo or a rename without updating the exporter produces a config that looks correct but fails with "unknown storage" at startup.

This rule fires when an exporter's `sending_queue.storage` references an extension name that does not exist in `extensions:`.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
extensions:
  file_storage:
    directory: /var/lib/otelcol/queue

exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      storage: file_storage/queue         # not defined — real name is "file_storage"
```

:::

:::tip[Correct]

```yaml
extensions:
  file_storage/queue:
    directory: /var/lib/otelcol/queue

exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      storage: file_storage/queue

service:
  extensions: [file_storage/queue]
```

:::

## When Not To Use It

Never — the Collector will refuse to start.

## Related Rules

- [OTEL-048](../exporter/otel-048) — `sending_queue` explicitly disabled
- [OTEL-062](./otel-062) — extension in `service.extensions` but not defined
- [OTEL-065](../reliability/otel-065) — `sending_queue` without persistent storage

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — file_storage extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/storage/filestorage/README.md)

## Resources

- Rule source: [`policy/main/extension.rego`](https://github.com/starkross/augur/blob/main/policy/main/extension.rego)
