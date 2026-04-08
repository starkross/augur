---
id: otel-034
title: "OTEL-034: CORS allowed_origins contains wildcard *"
sidebar_label: OTEL-034
description: A wildcard CORS origin lets any browser on any site push telemetry to this Collector.
---

# OTEL-034: CORS `allowed_origins` contains wildcard `*`

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Setting CORS `allowed_origins: ["*"]` on an OTLP/HTTP receiver means any web page anywhere can make cross-origin requests to the Collector. Unless the Collector is deliberately a public telemetry ingest, that is a full CSRF surface: any browser running any page can push arbitrary spans, metrics, or logs into your pipelines. List your actual front-end origins explicitly.

This rule fires when a receiver protocol has `cors.allowed_origins` containing the literal string `*`.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: "0.0.0.0:4318"
        cors:
          allowed_origins: ["*"]
```

:::

:::tip[Correct]

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: "0.0.0.0:4318"
        cors:
          allowed_origins:
            - https://app.example.com
            - https://admin.example.com
```

:::

## When Not To Use It

A genuinely public, unauthenticated ingest endpoint (rare). In that case pair the wildcard with a rate limiter and a strict schema so abuse is bounded.

## Related Rules

- [OTEL-010](../core/otel-010) — receivers should not bind to `0.0.0.0`
- [OTEL-033](./otel-033) — receiver on non-localhost endpoint without TLS

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP receiver](https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/otlpreceiver/README.md)
- [MDN — CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
