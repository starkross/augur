---
id: otel-051
title: "OTEL-051: sending_queue.num_consumers below 2"
sidebar_label: OTEL-051
description: A single consumer blocks the whole queue on one slow request.
---

# OTEL-051: `sending_queue.num_consumers` below 2

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`num_consumers` controls how many worker goroutines drain the sending queue in parallel. With just one worker, a single slow export — a tail-latency spike from the backend, a TLS handshake under load — stalls the queue until it completes. Two is the minimum for real concurrency; production workloads usually want 4–10.

This rule fires when `sending_queue.num_consumers < 2`.

## Options

| Field | Constraint |
|------|------------|
| `sending_queue.num_consumers` | Should be ≥ 2 |

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      num_consumers: 1
      queue_size: 5000
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 5000
```

:::

## When Not To Use It

Backends that serialize all writes and fall over under concurrency. In that case you are trading throughput for correctness — acknowledge the warning only for that specific exporter.

## Related Rules

- [OTEL-048](./otel-048) — `sending_queue` explicitly disabled
- [OTEL-049](./otel-049) — `sending_queue.queue_size` below 10
- [OTEL-050](./otel-050) — `sending_queue.queue_size` above 50000

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
