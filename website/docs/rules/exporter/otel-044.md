---
id: otel-044
title: "OTEL-044: OTLP gRPC exporter endpoint has http(s):// scheme"
sidebar_label: OTEL-044
description: The OTLP gRPC exporter wants a bare host:port — no URL scheme.
---

# OTEL-044: OTLP gRPC exporter endpoint has `http(s)://` scheme

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

The `otlp` exporter speaks gRPC, and its `endpoint` field expects a bare `host:port`. If you include an `http://` or `https://` prefix, the exporter's gRPC dialer fails at startup with a confusing `invalid target` error and the pipeline never starts. Use the `otlphttp` exporter if you need a URL scheme.

This rule fires when an `otlp` exporter's `endpoint` starts with `http://` or `https://`.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
exporters:
  otlp/backend:
    endpoint: https://backend.example.com:4317   # scheme not allowed
```

:::

:::tip[Correct]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
    tls:
      insecure: false
```

:::

## When Not To Use It

Never — this will fail at Collector start-up.

## Related Rules

- [OTEL-045](./otel-045) — OTLP gRPC endpoint missing port number
- [OTEL-046](./otel-046) — OTLP HTTP endpoint missing URL scheme
- [OTEL-047](./otel-047) — OTLP HTTP exporter using gRPC port 4317

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlpexporter/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
