---
id: otel-013
title: "OTEL-013: batch processor should be last in pipeline"
sidebar_label: OTEL-013
description: Putting batch last means every other processor sees un-batched data and batches reflect the final shape.
---

# OTEL-013: `batch` processor should be last in pipeline

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `batch` processor groups items right before they are exported. If it runs *before* a filter, transform, or sampler, the downstream processors see bigger objects than they need, do more work, and can fragment batches in ways the exporter was not sized for. The conventional ordering is `memory_limiter → …everything else… → batch`.

This rule fires when a pipeline contains `batch` and it is not the last processor.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, attributes]  # batch not last
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

Advanced setups with connectors that re-route data mid-pipeline may legitimately need a second batching stage. In that case disable the rule for the affected pipeline only.

## Related Rules

- [OTEL-014](./otel-014) — `memory_limiter` should be first processor
- [OTEL-038](../pipeline/otel-038) — filter processor after batch
- [OTEL-039](../pipeline/otel-039) — transform/attributes processor after batch
- [OTEL-043](../pipeline/otel-043) — batch before `tail_sampling`/`groupbytrace`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — batch processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/batchprocessor/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
