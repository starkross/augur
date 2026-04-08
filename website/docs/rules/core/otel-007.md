---
id: otel-007
title: "OTEL-007: Every pipeline must have receivers and exporters"
sidebar_label: OTEL-007
description: A pipeline with no receivers or no exporters transmits nothing.
---

# OTEL-007: Every pipeline must have receivers and exporters

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

A pipeline is only useful if data can enter it (at least one receiver) and leave it (at least one exporter). A pipeline missing either side is a dead end — it costs memory and config review time while sending nothing downstream. Processors are optional; receivers and exporters are not.

This rule fires when a pipeline in `service.pipelines` has zero receivers or zero exporters.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: []               # no exporters
    metrics:
      receivers: []               # no receivers
      exporters: [otlp/backend]
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
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Never. If a pipeline has no receivers or no exporters, delete it — the config is easier to read without dead code.

## Related Rules

- [OTEL-006](./otel-006) — `service.pipelines` must be defined
- [OTEL-020](./otel-020) — unused receiver
- [OTEL-021](./otel-021) — unused exporter

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — pipelines](https://opentelemetry.io/docs/collector/configuration/#service)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
