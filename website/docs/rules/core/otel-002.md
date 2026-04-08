---
id: otel-002
title: "OTEL-002: memory_limiter must be included in every pipeline"
sidebar_label: OTEL-002
description: Declaring memory_limiter is not enough — each pipeline must reference it.
---

# OTEL-002: memory_limiter must be included in every pipeline

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Declaring `memory_limiter` under `processors` only defines it — it does nothing until a pipeline actually includes it in its `processors:` list. A pipeline that skips `memory_limiter` has no back-pressure, so it keeps accepting data even when the Collector is already over its memory budget.

This rule fires when the `memory_limiter` processor is declared but at least one pipeline in `service.pipelines` does not include it.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]            # missing memory_limiter
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
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

There is no meaningful exception. If you are running multiple pipelines, every one of them should push back under memory pressure — otherwise the pipelines that skip the limiter will keep consuming and starve the ones that honor it.

## Related Rules

- [OTEL-001](./otel-001) — `memory_limiter` processor must be configured
- [OTEL-014](./otel-014) — `memory_limiter` should be first processor
- [OTEL-022](./otel-022) — unused processor

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)
- [OpenTelemetry Collector — pipelines](https://opentelemetry.io/docs/collector/configuration/#service)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
