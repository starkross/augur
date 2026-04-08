---
id: otel-053
title: "OTEL-053: Retry max_elapsed_time set to 0 (infinite retries)"
sidebar_label: OTEL-053
description: Infinite retries guarantee unbounded queue growth during an outage.
---

# OTEL-053: Retry `max_elapsed_time` set to 0 (infinite retries)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`retry_on_failure.max_elapsed_time: 0` disables the retry cutoff — the exporter will keep retrying the same batch forever. During a real backend outage that means the sending queue fills with stuck batches, the exporter holds them indefinitely, and either `memory_limiter` starts refusing new work or the process OOMs. Set a finite upper bound (5–10 minutes is typical) so the exporter eventually gives up and lets the queue drain.

This rule fires when `retry_on_failure.max_elapsed_time` is a zero duration.

## Options

| Field | Constraint |
|------|------------|
| `retry_on_failure.max_elapsed_time` | Must be > 0 |

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    retry_on_failure:
      enabled: true
      max_elapsed_time: 0s         # never gives up
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 300s
```

:::

## When Not To Use It

Never — always set a finite cutoff.

## Related Rules

- [OTEL-017](../core/otel-017) — exporter missing `retry_on_failure`/`sending_queue`
- [OTEL-050](./otel-050) — `sending_queue.queue_size` above 50000

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
