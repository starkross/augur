---
id: otel-052
title: "OTEL-052: Compression disabled for network exporter"
sidebar_label: OTEL-052
description: gzip typically reduces OTLP bandwidth by 70–90% — leaving it off is expensive.
---

# OTEL-052: Compression disabled for network exporter

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Telemetry payloads are highly repetitive (resource attributes, span names, metric labels) and compress extremely well — gzip commonly reduces OTLP batches by 70–90%. Setting `compression: none` on a network exporter throws all of that away, inflating egress cost and slowing every request proportionally. Pull-based exporters (`debug`, `logging`, `prometheus`, `prometheusremotewrite`, `file`) are exempt.

This rule fires when a non-pull-based exporter sets `compression: none`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    compression: none
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    compression: gzip
```

:::

## When Not To Use It

Network paths with hardware offloaded compression (rare), or intentionally disabled compression during a CPU profiling exercise.

## Related Rules

- [OTEL-017](../core/otel-017) — exporter missing `retry_on_failure`/`sending_queue`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlpexporter/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
