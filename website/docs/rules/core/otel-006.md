---
id: otel-006
title: "OTEL-006: service.pipelines must be defined"
sidebar_label: OTEL-006
description: A Collector config with no pipelines does nothing.
---

# OTEL-006: `service.pipelines` must be defined

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

`service.pipelines` is what actually wires receivers, processors, and exporters together. A config without any pipelines is syntactically valid but operationally inert — the Collector starts, logs nothing useful, and silently drops every signal. This is almost always an unfinished or broken config.

This rule fires when `service.pipelines` is missing entirely.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

exporters:
  otlp/backend:
    endpoint: backend:4317

service:
  telemetry:
    logs:
      level: info
  # pipelines block is missing
```

:::

:::tip[Correct]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Never. If you do not want any telemetry to flow, do not run a Collector at all.

## Related Rules

- [OTEL-007](./otel-007) — every pipeline must have receivers and exporters
- [OTEL-020](./otel-020) — unused receiver
- [OTEL-021](./otel-021) — unused exporter
- [OTEL-022](./otel-022) — unused processor

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — service configuration](https://opentelemetry.io/docs/collector/configuration/#service)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
