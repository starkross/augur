---
id: otel-024
title: "OTEL-024: batch send_batch_max_size < send_batch_size"
sidebar_label: OTEL-024
description: The hard cap must be at least as large as the soft trigger.
---

# OTEL-024: `batch` `send_batch_max_size` < `send_batch_size`

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

`send_batch_size` is the soft trigger (flush when the batch reaches this size) and `send_batch_max_size` is the hard cap (split anything larger). If the cap is smaller than the trigger, the trigger can never fire on a full batch — every batch is capped early, flushed at a size lower than intended, and the processor's throughput logic is effectively broken. The cap must be ≥ the trigger.

This rule fires when `send_batch_max_size < send_batch_size` on any `batch` processor.

## Options

| Field | Constraint |
|------|------------|
| `send_batch_max_size` | Must be ≥ `send_batch_size` |

## Examples

:::danger[Incorrect]

```yaml
processors:
  batch:
    timeout: 1s
    send_batch_size: 8192
    send_batch_max_size: 1024       # smaller than send_batch_size
```

:::

:::tip[Correct]

```yaml
processors:
  batch:
    timeout: 1s
    send_batch_size: 8192
    send_batch_max_size: 16384
```

:::

## When Not To Use It

Never — this configuration cannot behave correctly.

## Related Rules

- [OTEL-023](./otel-023) — `batch` `send_batch_max_size` unset
- [OTEL-003](../core/otel-003) — `batch` processor must be configured
- [OTEL-013](../core/otel-013) — `batch` processor should be last in pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
