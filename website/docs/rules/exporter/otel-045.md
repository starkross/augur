---
id: otel-045
title: "OTEL-045: OTLP gRPC endpoint missing port number"
sidebar_label: OTEL-045
description: Without a port, the exporter connects to whatever default it picks.
---

# OTEL-045: OTLP gRPC endpoint missing port number

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `otlp` exporter's `endpoint` is a `host:port` string. Leaving the port off means the gRPC dialer either fails or silently connects to a default port you did not intend, and the failure mode is different depending on the Go version. Always list the port explicitly — 4317 for standard OTLP gRPC.

This rule fires when an `otlp` exporter has an `endpoint` that does not contain `:`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com        # no port
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
```

:::

## When Not To Use It

Never — always specify the port.

## Related Rules

- [OTEL-044](./otel-044) — OTLP gRPC exporter endpoint has `http(s)://` scheme
- [OTEL-046](./otel-046) — OTLP HTTP endpoint missing URL scheme
- [OTEL-047](./otel-047) — OTLP HTTP exporter using gRPC port 4317

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlpexporter/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
