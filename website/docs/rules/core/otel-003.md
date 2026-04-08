---
id: otel-003
title: "OTEL-003: batch processor must be configured"
sidebar_label: OTEL-003
description: The batch processor is required for efficient, cost-effective export.
---

# OTEL-003: `batch` processor must be configured

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

The `batch` processor groups spans, metrics, and logs before they leave the Collector. Without it, each signal is exported individually — which explodes network round-trips, multiplies backend ingestion cost, and typically triggers rate-limiting on the downstream system. Every production config should have a `batch` processor.

This rule fires when no `batch` processor is declared in the top-level `processors` section.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
```

:::

:::tip[Correct]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s
    send_batch_size: 1024
```

:::

## When Not To Use It

Low-volume development environments where each trace or metric round-trip is easier to debug un-batched. For production, keep batching on.

## Related Rules

- [OTEL-013](./otel-013) — `batch` processor should be last in pipeline
- [OTEL-023](../memory/otel-023) — `batch` `send_batch_max_size` unset
- [OTEL-024](../memory/otel-024) — `batch` `send_batch_max_size` < `send_batch_size`
- [OTEL-025](../memory/otel-025) — `batch` timeout below 100ms
- [OTEL-026](../memory/otel-026) — `batch` timeout above 60s

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
