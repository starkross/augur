---
id: otel-049
title: "OTEL-049: sending_queue.queue_size below 10"
sidebar_label: OTEL-049
description: A queue that holds fewer than 10 items fills up before the first retry completes.
---

# OTEL-049: `sending_queue.queue_size` below 10

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The sending queue is sized in number of batches, not items — so 10 is already a small number. Going below 10 means one slow export keeps the queue full, back-pressure kicks in immediately, and the Collector drops data that would otherwise have been retried. Tune based on your target backend's recovery time, but do not go below 10.

This rule fires when `sending_queue.queue_size < 10`.

## Options

| Field | Constraint |
|------|------------|
| `sending_queue.queue_size` | Should be ≥ 10 |

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 5            # too small
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      queue_size: 5000
      num_consumers: 10
```

:::

## When Not To Use It

A tiny development exporter where data loss on failure is fine.

## Related Rules

- [OTEL-048](./otel-048) — `sending_queue` explicitly disabled
- [OTEL-050](./otel-050) — `sending_queue.queue_size` above 50000
- [OTEL-051](./otel-051) — `sending_queue.num_consumers` below 2

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
