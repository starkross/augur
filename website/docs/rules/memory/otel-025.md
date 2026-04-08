---
id: otel-025
title: "OTEL-025: batch timeout below 100ms"
sidebar_label: OTEL-025
description: Very short timeouts flush before any meaningful batching happens.
---

# OTEL-025: `batch` timeout below 100ms

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`timeout` is the maximum time a batch can wait before being flushed. Anything under ~100ms means most batches flush on the timeout rather than on `send_batch_size`, so you pay the per-batch overhead (gRPC headers, compression setup, TLS handshakes for keepalive-less exporters) on tiny batches. If you need low end-to-end latency, set other knobs; do not starve the batcher.

This rule fires when `batch.timeout < 100ms`.

## Options

| Field | Constraint |
|------|------------|
| `timeout` | Should be ≥ 100ms |

## Examples

:::warning[Avoid]

```yaml
processors:
  batch:
    timeout: 10ms                    # below 100ms
    send_batch_size: 1024
```

:::

:::tip[Prefer]

```yaml
processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
```

:::

## When Not To Use It

Ultra-low-latency tracing pipelines where you deliberately trade throughput for freshness and have measured that a sub-100ms timeout is actually needed.

## Related Rules

- [OTEL-023](./otel-023) — `batch` `send_batch_max_size` unset
- [OTEL-024](./otel-024) — `send_batch_max_size` < `send_batch_size`
- [OTEL-026](./otel-026) — `batch` timeout above 60s

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
