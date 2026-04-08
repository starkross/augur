---
id: otel-050
title: "OTEL-050: sending_queue.queue_size above 50000 (OOM risk)"
sidebar_label: OTEL-050
description: A very large queue can outlive the memory_limiter and OOM the Collector.
---

# OTEL-050: `sending_queue.queue_size` above 50000 (OOM risk)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The sending queue holds batches in memory. At 50,000 batches × realistic batch sizes you are looking at gigabytes of buffered data that `memory_limiter` cannot see (the limiter protects the processor side, not the exporter side). When the downstream is slow for long enough, the queue fills, then the process OOMs. Use a persistent storage extension ([OTEL-065](../reliability/otel-065)) or cap the queue.

This rule fires when `sending_queue.queue_size > 50000`.

## Options

| Field | Constraint |
|------|------------|
| `sending_queue.queue_size` | Should be ≤ 50000 |

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 100000         # oversized
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 10000
      num_consumers: 10
      storage: file_storage/queue   # persistent spill-over
```

:::

## When Not To Use It

An exporter writing to a very slow cold sink (object storage) where you genuinely need a large in-memory buffer. Pair with persistent storage.

## Related Rules

- [OTEL-049](./otel-049) — `sending_queue.queue_size` below 10
- [OTEL-048](./otel-048) — `sending_queue` explicitly disabled
- [OTEL-065](../reliability/otel-065) — `sending_queue` without persistent storage

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
