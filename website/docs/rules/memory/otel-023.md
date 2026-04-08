---
id: otel-023
title: "OTEL-023: batch send_batch_max_size unset (unlimited)"
sidebar_label: OTEL-023
description: Without an upper bound, a single batch can outgrow exporter limits.
---

# OTEL-023: `batch` `send_batch_max_size` unset (unlimited)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`send_batch_size` is the soft trigger that fires a flush; `send_batch_max_size` is the hard cap that splits oversized batches. With no cap, one very large burst (or a slow exporter that lets batches accumulate) can produce a single payload that the downstream refuses (gRPC `ResourceExhausted`) or that OOMs the exporter side. A good default is roughly twice `send_batch_size`.

This rule fires when a `batch` processor has no `send_batch_max_size` set.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
processors:
  batch:
    timeout: 1s
    send_batch_size: 8192
    # no send_batch_max_size
```

:::

:::tip[Prefer]

```yaml
processors:
  batch:
    timeout: 1s
    send_batch_size: 8192
    send_batch_max_size: 16384
```

:::

## When Not To Use It

Environments where the downstream exporter is known to accept arbitrarily large payloads and you would rather never split a batch. Rare in practice.

## Related Rules

- [OTEL-003](../core/otel-003) — `batch` processor must be configured
- [OTEL-024](./otel-024) — `send_batch_max_size` < `send_batch_size`
- [OTEL-025](./otel-025) — `batch` timeout below 100ms
- [OTEL-026](./otel-026) — `batch` timeout above 60s

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
