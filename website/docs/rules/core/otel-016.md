---
id: otel-016
title: "OTEL-016: Telemetry log level set to debug"
sidebar_label: OTEL-016
description: Debug-level logs in production are expensive and noisy.
---

# OTEL-016: Telemetry log level set to `debug`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`service.telemetry.logs.level: debug` turns on high-volume internal logging from every Collector component. In production this burns CPU on string formatting, multiplies log storage cost, and makes real problems harder to find. `info` is the right default for production; `debug` is for short-lived investigations.

This rule fires when `service.telemetry.logs.level == "debug"`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  telemetry:
    logs:
      level: debug
```

:::

:::tip[Prefer]

```yaml
service:
  telemetry:
    logs:
      level: info
```

:::

## When Not To Use It

Short debugging sessions in non-production environments. Set the level back to `info` before deploying.

## Related Rules

- [OTEL-015](./otel-015) — `debug`/`logging` exporter detected
- [OTEL-069](../lifecycle/otel-069) — telemetry metrics level set to `none`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — internal telemetry](https://opentelemetry.io/docs/collector/internal-telemetry/)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
