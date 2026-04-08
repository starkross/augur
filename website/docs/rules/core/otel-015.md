---
id: otel-015
title: "OTEL-015: debug/logging exporter detected"
sidebar_label: OTEL-015
description: The debug exporter is a development tool — it should not ship to production.
---

# OTEL-015: `debug`/`logging` exporter detected

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `debug` exporter (formerly `logging`) prints telemetry to stdout. It is invaluable during local development but in production it floods logs, pins CPU on serialization, and fills disks. A Collector with `debug` wired into a live pipeline is almost always a debugging leftover that should be removed.

This rule fires when an exporter named `debug` or `logging` is present in the `exporters` block.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  debug: {}
  otlp/backend:
    endpoint: backend:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug, otlp/backend]
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

Local development and ad-hoc debugging sessions. Use a separate config file for that case instead of shipping `debug` to production.

## Related Rules

- [OTEL-021](./otel-021) — unused exporter
- [OTEL-071](../lifecycle/otel-071) — `logging` exporter deprecated (renamed to `debug`)

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — debug exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/debugexporter/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
