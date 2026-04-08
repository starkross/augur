---
id: otel-032
title: "OTEL-032: insecure_skip_verify enabled"
sidebar_label: OTEL-032
description: Skipping certificate verification makes TLS meaningless.
---

# OTEL-032: `insecure_skip_verify` enabled

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`insecure_skip_verify: true` tells the Collector to accept any certificate from the peer — expired, self-signed, for the wrong hostname, or issued by a CA you do not trust. That silently converts "mTLS to our backend" into "encrypted connection to anyone who can intercept the traffic." Use a proper CA bundle or pin a specific cert instead.

This rule fires on any exporter or receiver protocol that sets `tls.insecure_skip_verify: true`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
    tls:
      insecure_skip_verify: true    # skips verification entirely
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
    tls:
      ca_file: /etc/ssl/certs/backend-ca.pem
```

:::

## When Not To Use It

One-off local development or a smoke test inside a throwaway environment. Never in production.

## Related Rules

- [OTEL-031](./otel-031) — TLS `min_version` below 1.2
- [OTEL-033](./otel-033) — receiver on non-localhost endpoint without TLS
- [OTEL-018](../core/otel-018) — OTLP exporter without TLS on non-local endpoint

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — TLS configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
