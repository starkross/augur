---
id: otel-021
title: "OTEL-021: Unused exporter"
sidebar_label: OTEL-021
description: An exporter not referenced in any pipeline sends nothing.
---

# OTEL-021: Unused exporter

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

An exporter defined but not referenced in any pipeline never flushes data. The Collector allocates it, the CI passes, and nobody notices that the new backend is receiving zero telemetry. Delete or wire the exporter up.

This rule fires when an exporter is declared but not used by any pipeline.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
  otlp/unused:                         # never wired into a pipeline
    endpoint: other.example.com:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Pre-staged exporters for a cutover (old backend → new backend) may sit unused for a window. Acknowledge the warning until the cutover lands.

## Related Rules

- [OTEL-020](./otel-020) — unused receiver
- [OTEL-022](./otel-022) — unused processor

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporters](https://opentelemetry.io/docs/collector/configuration/#exporters)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
