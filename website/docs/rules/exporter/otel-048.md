---
id: otel-048
title: "OTEL-048: sending_queue explicitly disabled"
sidebar_label: OTEL-048
description: Disabling the sending queue drops data on the first transient failure.
---

# OTEL-048: `sending_queue` explicitly disabled

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`sending_queue` is the exporter helper's in-memory buffer that absorbs bursts and retries. Setting `enabled: false` turns it off — any request that fails once is dropped on the floor. That is a dramatic downgrade from the default behavior and usually a leftover from debugging. Pull-based exporters (`debug`, `logging`, `prometheus`, `prometheusremotewrite`, `file`) do not use a sending queue and are exempt.

This rule fires when a non-pull-based exporter has `sending_queue.enabled: false`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    sending_queue:
      enabled: false          # disabled
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

Never for production. If you need synchronous behavior for a test, use a local file exporter or `debug` instead of turning off the queue on a real network exporter.

## Related Rules

- [OTEL-017](../core/otel-017) — exporter missing `retry_on_failure`/`sending_queue`
- [OTEL-049](./otel-049) — `sending_queue.queue_size` below 10
- [OTEL-050](./otel-050) — `sending_queue.queue_size` above 50000
- [OTEL-051](./otel-051) — `sending_queue.num_consumers` below 2
- [OTEL-065](../reliability/otel-065) — `sending_queue` without persistent storage

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
