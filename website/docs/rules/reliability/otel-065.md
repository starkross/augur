---
id: otel-065
title: "OTEL-065: sending_queue without persistent storage"
sidebar_label: OTEL-065
description: An in-memory queue loses everything on restart or crash.
---

# OTEL-065: `sending_queue` without persistent storage

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

With `sending_queue` enabled but no `storage` extension wired in, the queue lives entirely in RAM. Every restart — rolling deploy, OOM, node reschedule, cgroup kill — drops every batch that had not yet been exported. For any pipeline that carries transactional or billing-grade data, back the queue with the `file_storage` extension so unsent batches survive a restart.

This rule fires when an exporter has `sending_queue.enabled != false` and no `sending_queue.storage` configured.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 5000
      # no storage
```

:::

:::tip[Prefer]

```yaml
extensions:
  file_storage/queue:
    directory: /var/lib/otelcol/queue

exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 5000
      storage: file_storage/queue

service:
  extensions: [file_storage/queue]
```

:::

## When Not To Use It

Short-lived non-durable telemetry (local debug runs, ephemeral CI pipelines) where losing data on restart is acceptable.

## Related Rules

- [OTEL-017](../core/otel-017) — exporter missing `retry_on_failure`/`sending_queue`
- [OTEL-048](../exporter/otel-048) — `sending_queue` explicitly disabled
- [OTEL-050](../exporter/otel-050) — `sending_queue.queue_size` above 50000
- [OTEL-066](../extension/otel-066) — `sending_queue.storage` references undefined extension

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — file_storage extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/storage/filestorage/README.md)
- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/reliability.rego`](https://github.com/starkross/augur/blob/main/policy/main/reliability.rego)
