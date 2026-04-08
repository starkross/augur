---
id: otel-063
title: "OTEL-063: tail_sampling without groupbytrace"
sidebar_label: OTEL-063
description: tail_sampling needs complete traces — groupbytrace is what assembles them.
---

# OTEL-063: `tail_sampling` without `groupbytrace`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`tail_sampling` makes sampling decisions based on the *whole* trace — error spans, latency outliers, root span attributes. It can only do that when it sees every span for a trace in the same batch. `groupbytrace` is the processor responsible for buffering spans by trace ID and handing complete traces over to the sampler. Running `tail_sampling` alone means it sees partial traces and makes inconsistent decisions — some spans sampled in, some dropped, same trace.

This rule fires when a pipeline contains `tail_sampling` but not `groupbytrace`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, tail_sampling, batch]   # no groupbytrace
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

A head-sampling-only pipeline that happens to include a degenerate `tail_sampling` policy (rare). In practice, if you use `tail_sampling` you need `groupbytrace`.

## Related Rules

- [OTEL-043](../pipeline/otel-043) — batch before `tail_sampling`/`groupbytrace`
- [OTEL-064](./otel-064) — both `probabilistic_sampler` and `tail_sampling` in same pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — tail_sampling processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/tailsamplingprocessor/README.md)
- [OpenTelemetry Collector Contrib — groupbytrace processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/groupbytraceprocessor/README.md)

## Resources

- Rule source: [`policy/main/reliability.rego`](https://github.com/starkross/augur/blob/main/policy/main/reliability.rego)
