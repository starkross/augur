---
id: otel-040
title: "OTEL-040: Circular pipeline dependency via connectors"
sidebar_label: OTEL-040
description: Connector cycles make telemetry loop forever.
---

# OTEL-040: Circular pipeline dependency via connectors

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Connectors turn one pipeline's exporter side into another pipeline's receiver side. That is a directed graph — and a cycle in that graph means telemetry loops indefinitely between pipelines, consuming CPU and memory until one of the other rules (memory_limiter, queue size) terminates the loop. augur walks the pipeline graph and blocks any config whose connector edges form a cycle.

This rule fires when the connector-pipeline graph contains at least one cycle.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
connectors:
  forward/a_to_b: {}
  forward/b_to_a: {}

service:
  pipelines:
    traces/a:
      receivers: [otlp]
      exporters: [forward/a_to_b]
    traces/b:
      receivers: [forward/a_to_b]
      exporters: [forward/b_to_a]
    traces/c:
      receivers: [forward/b_to_a]
      exporters: [forward/a_to_b]   # cycle: a_to_b -> b_to_a -> a_to_b
```

:::

:::tip[Correct]

```yaml
connectors:
  forward/a_to_b: {}

service:
  pipelines:
    traces/a:
      receivers: [otlp]
      exporters: [forward/a_to_b]
    traces/b:
      receivers: [forward/a_to_b]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Never — a connector cycle is always a bug.

## Related Rules

- [OTEL-041](./otel-041) — routing connector without `default_pipelines`
- [OTEL-042](./otel-042) — duplicate processor in same pipeline

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — connectors](https://opentelemetry.io/docs/collector/configuration/#connectors)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
