---
id: otel-058
title: "OTEL-058: Multiple receivers bound to the same endpoint"
sidebar_label: OTEL-058
description: Two receivers competing for the same port will crash the Collector at startup.
---

# OTEL-058: Multiple receivers bound to the same endpoint

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Two receivers pointed at the same `host:port` will both try to `bind(2)` it at Collector startup. One succeeds, one fails with `address already in use`, and the Collector exits. This is usually the result of a merge that duplicated an OTLP block, or a template that emits the same receiver twice under different names.

This rule fires when two distinct receivers have the same endpoint string.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
  otlp/secondary:
    protocols:
      grpc:
        endpoint: localhost:4317       # collides
```

:::

:::tip[Correct]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
  otlp/secondary:
    protocols:
      grpc:
        endpoint: localhost:14317
```

:::

## When Not To Use It

Never — the config cannot start.

## Related Rules

- [OTEL-010](../core/otel-010) — receivers should not bind to `0.0.0.0`
- [OTEL-020](../core/otel-020) — unused receiver

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — receivers](https://opentelemetry.io/docs/collector/configuration/#receivers)

## Resources

- Rule source: [`policy/main/receiver.rego`](https://github.com/starkross/augur/blob/main/policy/main/receiver.rego)
