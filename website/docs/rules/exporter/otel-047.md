---
id: otel-047
title: "OTEL-047: OTLP HTTP exporter using gRPC port 4317"
sidebar_label: OTEL-047
description: Port 4317 is for gRPC, port 4318 is for HTTP.
---

# OTEL-047: OTLP HTTP exporter using gRPC port 4317 (HTTP is 4318)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

By convention (and per the OTLP spec) port 4317 is OTLP/gRPC and port 4318 is OTLP/HTTP. Pointing `otlphttp` at `:4317` almost always means the author mixed up the two exporters — the TCP connection will succeed, but the HTTP request will return protocol errors from the gRPC server on the other side.

This rule fires when an `otlphttp` exporter's `endpoint` ends with `:4317` or contains `:4317/`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlphttp/backend:
    endpoint: https://backend.example.com:4317    # gRPC port
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlphttp/backend:
    endpoint: https://backend.example.com:4318
```

:::

## When Not To Use It

A backend that genuinely serves OTLP/HTTP on port 4317 — rare and non-standard, but possible. In that case disable the rule for that exporter.

## Related Rules

- [OTEL-044](./otel-044) — OTLP gRPC exporter endpoint has `http(s)://` scheme
- [OTEL-045](./otel-045) — OTLP gRPC endpoint missing port number
- [OTEL-046](./otel-046) — OTLP HTTP endpoint missing URL scheme

## Version

Available since augur v0.1.0.

## Further Reading

- [OTLP specification — ports](https://opentelemetry.io/docs/specs/otlp/)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
