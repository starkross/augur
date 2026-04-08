---
id: otel-010
title: "OTEL-010: Receivers should not bind to 0.0.0.0"
sidebar_label: OTEL-010
description: Binding to 0.0.0.0 exposes receivers on every interface.
---

# OTEL-010: Receivers should not bind to `0.0.0.0`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Binding a receiver to `0.0.0.0` tells the Collector to accept traffic on every network interface the host has — including public ones. For sidecars and single-node agents you almost always want `localhost` (or `127.0.0.1`) so only workloads on the same network namespace can push data. For gateway Collectors, bind to a specific private interface rather than `0.0.0.0`.

This rule fires when any string under a receiver contains the substring `0.0.0.0`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
      http:
        endpoint: "0.0.0.0:4318"
```

:::

:::tip[Prefer]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "localhost:4317"
      http:
        endpoint: "localhost:4318"
```

:::

## When Not To Use It

A gateway Collector that genuinely needs to accept traffic from arbitrary hosts on a trusted private network. In that case prefer binding to a specific interface (e.g. `10.0.0.5:4317`) over the wildcard, and combine with [OTEL-033](../security/otel-033) / [OTEL-031](../security/otel-031) to enforce TLS.

## Related Rules

- [OTEL-033](../security/otel-033) — receiver on non-localhost endpoint without TLS
- [OTEL-060](../extension/otel-060) — `zpages` endpoint bound to `0.0.0.0`
- [OTEL-070](../lifecycle/otel-070) — telemetry metrics address bound to `0.0.0.0`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — receivers](https://opentelemetry.io/docs/collector/configuration/#receivers)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
