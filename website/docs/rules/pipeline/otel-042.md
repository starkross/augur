---
id: otel-042
title: "OTEL-042: Duplicate processor in same pipeline"
sidebar_label: OTEL-042
description: Listing the same processor twice usually means a copy/paste mistake.
---

# OTEL-042: Duplicate processor in same pipeline

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Listing the same processor name twice in a pipeline's `processors:` list makes it run twice on the same data. For `memory_limiter` that is harmless; for anything that mutates (`attributes`, `transform`, `filter`, `tail_sampling`) it doubles side effects and is almost always a copy/paste artifact from a merge conflict or a refactor.

This rule fires when the same processor name appears more than once in a single pipeline's `processors:` list.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, attributes, attributes, batch]
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

A deliberate two-stage transform that legitimately runs the same processor instance twice — very rare. Use differently named instances (`attributes/first`, `attributes/second`) instead to make the intent explicit.

## Related Rules

- [OTEL-022](../core/otel-022) — unused processor
- [OTEL-013](../core/otel-013) — `batch` processor should be last in pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — pipelines](https://opentelemetry.io/docs/collector/configuration/#service)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
