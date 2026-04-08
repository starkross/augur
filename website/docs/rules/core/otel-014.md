---
id: otel-014
title: "OTEL-014: memory_limiter should be first processor in pipeline"
sidebar_label: OTEL-014
description: The memory limiter must be first so it can reject work before anything allocates.
---

# OTEL-014: `memory_limiter` should be first processor in pipeline

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The whole point of `memory_limiter` is to refuse new work when the Collector is already over its soft limit. If processors that allocate (transform, attributes, tail_sampling) run *before* the limiter, they have already spent memory by the time the limiter kicks in — defeating the back-pressure. Put `memory_limiter` first so it rejects early.

This rule fires when a pipeline contains `memory_limiter` and it is not the first processor.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [attributes, memory_limiter, batch]  # memory_limiter not first
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, attributes, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

There is no sane exception — if `memory_limiter` is in the pipeline at all, it should be first.

## Related Rules

- [OTEL-001](./otel-001) — `memory_limiter` processor must be configured
- [OTEL-002](./otel-002) — `memory_limiter` must be included in every pipeline
- [OTEL-013](./otel-013) — `batch` processor should be last in pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
