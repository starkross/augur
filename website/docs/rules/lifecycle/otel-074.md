---
id: otel-074
title: "OTEL-074: service.telemetry.metrics.address (deprecated, use readers config)"
sidebar_label: OTEL-074
description: The address field is deprecated in favor of the new readers configuration.
---

# OTEL-074: `service.telemetry.metrics.address` (deprecated, use `readers` config)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The old `service.telemetry.metrics.address` field is deprecated in favor of the newer `service.telemetry.metrics.readers` configuration, which supports multiple exporters (Prometheus, OTLP, file) and more expressive wiring. The old field still works but is scheduled for removal. Migrate to `readers` when you touch this block next.

This rule fires when `service.telemetry.metrics.address` is set.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  telemetry:
    metrics:
      level: detailed
      address: "localhost:8888"
```

:::

:::tip[Prefer]

```yaml
service:
  telemetry:
    metrics:
      level: detailed
      readers:
        - pull:
            exporter:
              prometheus:
                host: localhost
                port: 8888
```

:::

## When Not To Use It

Older Collector versions that do not yet understand `readers`. Check the release notes of your pinned version before migrating.

## Related Rules

- [OTEL-069](./otel-069) — telemetry metrics level set to `none`
- [OTEL-070](./otel-070) — telemetry metrics address bound to `0.0.0.0`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — internal telemetry](https://opentelemetry.io/docs/collector/internal-telemetry/)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
