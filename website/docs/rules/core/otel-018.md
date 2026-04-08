---
id: otel-018
title: "OTEL-018: OTLP exporter without TLS on non-local endpoint"
sidebar_label: OTEL-018
description: OTLP traffic to a remote backend should be encrypted.
---

# OTEL-018: OTLP exporter without TLS on non-local endpoint

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

OTLP carries spans, metrics, and logs that can contain request payloads, user identifiers, stack traces, and internal hostnames. Shipping any of that to a remote endpoint without TLS means passive network observers can read it and active ones can tamper with it. Local endpoints (`localhost`, `127.0.0.1`) are exempt.

This rule fires on `otlp*` exporters whose `endpoint` is not `localhost`/`127.0.0.1`, whose URL does not start with `https://`, and that have no `tls` block.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
    # no tls block
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend.example.com:4317
    tls:
      insecure: false
      ca_file: /etc/ssl/certs/backend-ca.pem
```

:::

## When Not To Use It

A Collector-to-Collector hop inside a trusted mesh that already provides mTLS via a sidecar proxy (e.g. Istio, Linkerd, Consul). In that case document the assumption explicitly and disable the rule only for the affected exporter.

## Related Rules

- [OTEL-031](../security/otel-031) — TLS `min_version` below 1.2
- [OTEL-032](../security/otel-032) — `insecure_skip_verify` enabled
- [OTEL-033](../security/otel-033) — receiver on non-localhost endpoint without TLS
- [OTEL-037](../security/otel-037) — inline `key_pem` detected

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — TLS configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
