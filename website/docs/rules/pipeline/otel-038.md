---
id: otel-038
title: "OTEL-038: Filter processor after batch (filter early to reduce waste)"
sidebar_label: OTEL-038
description: Filtering before batching means you do not pay to batch what you are about to drop.
---

# OTEL-038: Filter processor after batch

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Dropping a span, metric, or log after it has already been grouped into a batch wastes the work of batching, and in some configurations forces the batch to be split or re-packed. Run `filter` *before* `batch` so only the telemetry you intend to keep reaches the batcher.

This rule fires when a pipeline contains both `batch` and `filter` and the `filter` processor appears after `batch`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch, filter]   # filter after batch
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [memory_limiter, filter, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Rarely — you might want filtering late if the filter relies on an attribute that is only added after batching (unusual). In that case disable the rule for that specific pipeline and document why.

## Related Rules

- [OTEL-013](../core/otel-013) — `batch` processor should be last in pipeline
- [OTEL-039](./otel-039) — transform/attributes processor after batch
- [OTEL-043](./otel-043) — batch before `tail_sampling`/`groupbytrace`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — filter processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor/README.md)
- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
