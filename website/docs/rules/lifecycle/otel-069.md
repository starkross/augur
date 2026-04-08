---
id: otel-069
title: "OTEL-069: Telemetry metrics level set to none"
sidebar_label: OTEL-069
description: Disabling internal metrics means you cannot see the Collector's own health.
---

# OTEL-069: Telemetry metrics level set to `none`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The Collector's own metrics (`otelcol_process_cpu_seconds`, `otelcol_exporter_queue_size`, `otelcol_receiver_accepted_spans`, …) are what oncall looks at when the pipeline starts dropping data. Setting `service.telemetry.metrics.level: none` turns them all off, leaving operators blind. Keep it at `normal` (the default) or `detailed` for noisy environments.

This rule fires when `service.telemetry.metrics.level == "none"`.

## Options

| Field | Constraint |
|------|------------|
| `service.telemetry.metrics.level` | Should be `normal` or `detailed` |

## Examples

:::warning[Avoid]

```yaml
service:
  telemetry:
    metrics:
      level: none
```

:::

:::tip[Prefer]

```yaml
service:
  telemetry:
    metrics:
      level: normal
```

:::

## When Not To Use It

Never — the Collector's internal telemetry is essential for operating it.

## Related Rules

- [OTEL-016](../core/otel-016) — telemetry log level set to `debug`
- [OTEL-070](./otel-070) — telemetry metrics address bound to `0.0.0.0`
- [OTEL-074](./otel-074) — `service.telemetry.metrics.address` deprecated

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — internal telemetry](https://opentelemetry.io/docs/collector/internal-telemetry/)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
