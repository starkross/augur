---
id: otel-039
title: "OTEL-039: Transform/attributes processor after batch"
sidebar_label: OTEL-039
description: Rewriting attributes after batching is wasted work.
---

# OTEL-039: Transform/attributes processor after batch

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`transform` and `attributes` mutate individual items — they add resource attributes, rename fields, or redact PII. Running them *after* `batch` means the batcher already grouped items based on their pre-transform shape; the downstream exporter then has to unpack and re-serialize those items anyway. Run mutators before batching.

This rule fires when a pipeline contains `batch` and a `transform` or `attributes` processor appearing after it.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, attributes]
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

A transform that only runs on batch-level metadata (rare) can live after `batch`. In almost every other case, mutators belong before batching.

## Related Rules

- [OTEL-013](../core/otel-013) — `batch` processor should be last in pipeline
- [OTEL-038](./otel-038) — filter processor after batch
- [OTEL-043](./otel-043) — batch before `tail_sampling`/`groupbytrace`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — transform processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor/README.md)
- [OpenTelemetry Collector — attributes processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/attributesprocessor/README.md)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
