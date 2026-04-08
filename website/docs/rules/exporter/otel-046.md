---
id: otel-046
title: "OTEL-046: OTLP HTTP endpoint missing URL scheme"
sidebar_label: OTEL-046
description: The OTLP HTTP exporter needs a full URL (http:// or https://).
---

# OTEL-046: OTLP HTTP endpoint missing URL scheme

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`otlphttp` speaks HTTP, so its `endpoint` is a full URL. Passing a bare `host:port` makes the Collector guess a scheme — usually `http://`, which silently downgrades a production exporter to cleartext. Spell out `https://` (or `http://` for a local sink) explicitly.

This rule fires when an `otlphttp` exporter has an `endpoint` that does not start with `http://` or `https://` and is not a literal `${env:...}` reference.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlphttp/backend:
    endpoint: backend.example.com:4318    # missing scheme
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

Never — always include the scheme.

## Related Rules

- [OTEL-044](./otel-044) — OTLP gRPC exporter endpoint has `http(s)://` scheme
- [OTEL-045](./otel-045) — OTLP gRPC endpoint missing port number
- [OTEL-047](./otel-047) — OTLP HTTP exporter using gRPC port 4317

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP HTTP exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlphttpexporter/README.md)

## Resources

- Rule source: [`policy/main/exporter.rego`](https://github.com/starkross/augur/blob/main/policy/main/exporter.rego)
