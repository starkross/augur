---
id: otel-020
title: "OTEL-020: Unused receiver"
sidebar_label: OTEL-020
description: A receiver not referenced in any pipeline silently does nothing.
---

# OTEL-020: Unused receiver

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

A receiver defined in `receivers:` but not referenced in any `service.pipelines.*.receivers:` list will not receive anything — the Collector never starts it. These are almost always leftover from a refactor and cause a lot of confusion: "the config has `otlp`, why don't we see any traces?" Delete them or wire them up.

This rule fires when a receiver is declared but not used by any pipeline.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
  jaeger:                               # declared but unused
    protocols:
      grpc:

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Templated configs where a receiver is conditionally included by a deploy step may legitimately look "unused" to augur. If you can, gate the receiver at template level instead of leaving it dormant.

## Related Rules

- [OTEL-021](./otel-021) — unused exporter
- [OTEL-022](./otel-022) — unused processor
- [OTEL-007](./otel-007) — every pipeline must have receivers and exporters

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — receivers](https://opentelemetry.io/docs/collector/configuration/#receivers)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
