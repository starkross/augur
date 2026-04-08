---
id: otel-043
title: "OTEL-043: Batch before tail_sampling/groupbytrace"
sidebar_label: OTEL-043
description: Batching before trace grouping splits traces across batches.
---

# OTEL-043: Batch before `tail_sampling`/`groupbytrace`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`tail_sampling` and `groupbytrace` both need to see all spans for a trace in the same processor instance to make a correct decision. If `batch` runs first, it can put spans from the same trace into different batches — and spans from different batches arrive at the sampler at different times, sometimes after the sampler's decision window has already closed. Run trace-aware processors before `batch`.

This rule fires when a pipeline contains `batch` and a `tail_sampling` or `groupbytrace` processor appearing after it.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, tail_sampling]
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, groupbytrace, tail_sampling, batch]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Never for `tail_sampling` or `groupbytrace` — it will produce incorrect sampling decisions.

## Related Rules

- [OTEL-013](../core/otel-013) — `batch` processor should be last in pipeline
- [OTEL-038](./otel-038) — filter processor after batch
- [OTEL-039](./otel-039) — transform/attributes processor after batch
- [OTEL-063](../reliability/otel-063) — `tail_sampling` without `groupbytrace`
- [OTEL-064](../reliability/otel-064) — both `probabilistic_sampler` and `tail_sampling` in same pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — tail_sampling processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/tailsamplingprocessor/README.md)
- [OpenTelemetry Collector Contrib — groupbytrace processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/groupbytraceprocessor/README.md)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
