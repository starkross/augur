---
id: otel-022
title: "OTEL-022: Unused processor"
sidebar_label: OTEL-022
description: A processor not referenced in any pipeline does nothing.
---

# OTEL-022: Unused processor

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

A processor defined but not referenced in any pipeline is dead config — it is not applied to any data and no optimization happens. Most often this is the symptom of forgetting to wire `memory_limiter` or `batch` into a new pipeline you just added.

This rule fires when a processor is declared but not used by any pipeline.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s
  attributes:                          # declared but never used
    actions:
      - key: env
        value: production
        action: upsert

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s
  attributes:
    actions:
      - key: env
        value: production
        action: upsert

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, attributes, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Rarely. If a processor is staged for a future rollout, wire it into a disabled pipeline or gate it at templating time so intent is visible.

## Related Rules

- [OTEL-020](./otel-020) — unused receiver
- [OTEL-021](./otel-021) — unused exporter
- [OTEL-002](./otel-002) — `memory_limiter` must be included in every pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — processors](https://opentelemetry.io/docs/collector/configuration/#processors)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
