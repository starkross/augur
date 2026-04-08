---
id: otel-026
title: "OTEL-026: batch timeout above 60s"
sidebar_label: OTEL-026
description: Very long timeouts pin telemetry in memory and delay incident signals.
---

# OTEL-026: `batch` timeout above 60s

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

A `batch.timeout` above 60s means telemetry can sit in the Collector's in-memory queue for over a minute before it leaves. That delays alerts, inflates memory usage during traffic spikes, and multiplies data loss on crash because more data is buffered at any given moment. 60s is already generous for almost every workload.

This rule fires when `batch.timeout > 60s`.

## Options

| Field | Constraint |
|------|------------|
| `timeout` | Should be ≤ 60s |

## Examples

:::warning[Avoid]

```yaml
processors:
  batch:
    timeout: 120s                    # too long
```

:::

:::tip[Prefer]

```yaml
processors:
  batch:
    timeout: 5s
```

:::

## When Not To Use It

Cold-archive pipelines where you deliberately want huge batches written to object storage. Even then, 60s is usually enough.

## Related Rules

- [OTEL-025](./otel-025) — `batch` timeout below 100ms
- [OTEL-023](./otel-023) — `batch` `send_batch_max_size` unset

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
