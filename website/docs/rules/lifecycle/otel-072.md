---
id: otel-072
title: "OTEL-072: OpenCensus receiver/exporter (deprecated, migrate to OTLP)"
sidebar_label: OTEL-072
description: The OpenCensus protocol has been superseded by OTLP.
---

# OTEL-072: OpenCensus receiver/exporter (deprecated, migrate to OTLP)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

OpenCensus was the predecessor to OpenTelemetry; its wire protocol is in long-term deprecation. The `opencensus` receiver and exporter are still shipped for migration purposes but should not be used for new pipelines. Switch clients and backends to OTLP, which is the supported protocol for every current Collector component.

This rule fires when a receiver or exporter whose base name is `opencensus` is declared.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  opencensus:
    endpoint: localhost:55678

exporters:
  opencensus:
    endpoint: legacy-backend:55678
```

:::

:::tip[Prefer]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

exporters:
  otlp/backend:
    endpoint: backend:4317
```

:::

## When Not To Use It

Mid-migration Collectors that still need to bridge legacy OpenCensus clients to an OTLP backend. Once the clients are upgraded, remove the receiver/exporter.

## Related Rules

- [OTEL-071](./otel-071) — `logging` exporter deprecated
- [OTEL-073](./otel-073) — `memory_limiter` `ballast_size_mib` deprecated
- [OTEL-074](./otel-074) — `service.telemetry.metrics.address` deprecated

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry — migrating from OpenCensus](https://opentelemetry.io/docs/migration/opencensus/)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
