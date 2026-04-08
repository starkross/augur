---
id: otel-064
title: "OTEL-064: Both probabilistic_sampler and tail_sampling in same pipeline"
sidebar_label: OTEL-064
description: Stacking two samplers multiplies their rates and makes debugging impossible.
---

# OTEL-064: Both `probabilistic_sampler` and `tail_sampling` in same pipeline

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`probabilistic_sampler` is a head sampler (decision based on trace ID at ingestion); `tail_sampling` is a tail sampler (decision based on the whole trace). Running both means the effective sample rate is the product of the two, and the tail sampler can only see the traces that already survived the head sampler — so latency and error policies miss data. Pick one model for the pipeline and stick with it.

This rule fires when a pipeline contains both a `probabilistic_sampler` and a `tail_sampling` processor.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors:
        - memory_limiter
        - probabilistic_sampler
        - groupbytrace
        - tail_sampling
        - batch
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

Two-stage sampling pipelines (rare, advanced) where you have measured the combined rate and accepted the trace blind spot. Disable the rule for that pipeline and document the reasoning.

## Related Rules

- [OTEL-043](../pipeline/otel-043) — batch before `tail_sampling`/`groupbytrace`
- [OTEL-063](./otel-063) — `tail_sampling` without `groupbytrace`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — tail_sampling processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/tailsamplingprocessor/README.md)
- [OpenTelemetry Collector Contrib — probabilistic_sampler processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/probabilisticsamplerprocessor/README.md)

## Resources

- Rule source: [`policy/main/reliability.rego`](https://github.com/starkross/augur/blob/main/policy/main/reliability.rego)
